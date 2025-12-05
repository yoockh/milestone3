package service

import (
	"errors"
	"fmt"
	"log/slog"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	MinBidIncrement = 10000
	MaxRetries      = 3
)

var (
	mutex       sync.Map
	wibLocation *time.Location
)

func init() {
	var err error
	wibLocation, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// fallback to WIB
		wibLocation = time.FixedZone("WIB", 7*60*60)
	}
}

func getMutex(itemID int64) *sync.Mutex {
	m, _ := mutex.LoadOrStore(itemID, &sync.Mutex{})
	return m.(*sync.Mutex)
}

type bidService struct {
	redisRepo          repository.BidRedisRepository
	bidRepo            repository.BidRepository
	itemRepo           repository.AuctionItemRepository
	auctionSessionRepo repository.AuctionSessionRepository
	logger             *slog.Logger
}

type BidService interface {
	PlaceBid(sessionID, itemID, userID int64, amount float64, sessionEndTime time.Time) error
	GetHighestBid(sessionID, itemID int64) (float64, int64, error)

	SaveKeyToDB() error
	DeleteKeyValue() error
	CloseExpiredItemsWithoutBids() error
}

func NewBidService(r repository.BidRedisRepository, b repository.BidRepository, itemRepo repository.AuctionItemRepository, sessionRepo repository.AuctionSessionRepository, logger *slog.Logger) BidService {
	return &bidService{
		redisRepo:          r,
		bidRepo:            b,
		itemRepo:           itemRepo,
		auctionSessionRepo: sessionRepo,
		logger:             logger,
	}
}

func (s *bidService) PlaceBid(sessionID, itemID, userID int64, amount float64, sessionEndTime time.Time) error {
	if amount <= 0 {
		return ErrInvalidBidding
	}

	item, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		return ErrAuctionNotFound
	}

	if item.SessionID == nil || *item.SessionID != sessionID {
		return ErrInvalidAuction
	}

	if item.Status != "ongoing" {
		return ErrInvalidAuction
	}

	// validate session has started
	session, err := s.auctionSessionRepo.GetByID(sessionID)
	if err != nil {
		return ErrSessionNotFoundID
	}

	// Convert both to same timezone for comparison
	now := time.Now().In(wibLocation)

	// DB stores UTC, convert to WIB
	sessionStart := session.StartTime.In(wibLocation)
	sessionEnd := session.EndTime.In(wibLocation)

	if now.Before(sessionStart) {
		return ErrInvalidAuction
	}

	if now.After(sessionEnd) {
		return ErrInvalidAuction
	}

	if err = s.redisRepo.CheckDuplicateBid(userID, itemID, amount, 10*time.Second); err != nil {
		return ErrDuplicateBid
	}

	// lock per item
	mu := getMutex(itemID)
	mu.Lock()
	defer mu.Unlock()

	currentHighest, currentBid, err := s.redisRepo.GetHighestBid(sessionID, itemID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// validate amount > starting price
	if currentHighest == 0 {
		if amount < item.StartingPrice {
			return ErrBidTooLow
		}
	}

	if currentHighest > 0 && amount <= currentHighest {
		return ErrBidTooLow
	}

	if currentBid == userID {
		return ErrAlreadyHighestBidder
	}

	if currentHighest > 0 && amount < currentHighest+MinBidIncrement {
		return ErrBidTooLow
	}

	if err := s.redisRepo.SetHighestBid(sessionID, itemID, amount, userID, sessionEndTime); err != nil {
		s.logger.Error("failed to set highest bid", "error", err)
		return err
	}

	s.logger.Info("bid placed", "sessionID", sessionID, "itemID", itemID, "userID", userID, "amount", amount)

	return nil
}

func (s *bidService) GetHighestBid(sessionID, itemID int64) (float64, int64, error) {
	_, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		return 0, 0, ErrAuctionNotFound
	}
	return s.redisRepo.GetHighestBid(sessionID, itemID)
}

func parseKey(key string) (sessionID, itemID int64, err error) {
	parts := strings.Split(key, ":")

	if len(parts) != 5 {
		return 0, 0, errors.New("invalid key format")
	}

	if parts[0] != "active" || parts[1] != "auction" || parts[3] != "item" {
		return 0, 0, errors.New("invalid key format, expected: active:auction:{sessionID}:item:{itemID}")
	}

	sessionID, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid sessionID: %w", err)
	}

	itemID, err = strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid itemID: %w", err)
	}

	return sessionID, itemID, nil
}

func (s *bidService) SaveKeyToDB() error {
	pattern := "active:auction:*:item:*"

	keys, err := s.redisRepo.ScanKeys(pattern)
	if err != nil {
		s.logger.Error("failed to scan keys", "pattern", pattern, "error", err)
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	totalSavedAuctionItemToDB := 0
	for _, key := range keys {
		parsedSessionID, itemID, err := parseKey(key)
		if err != nil {
			continue
		}

		// fetch end_time from database for accuracy
		session, err := s.auctionSessionRepo.GetByID(parsedSessionID)
		if err != nil {
			continue
		}

		// Convert both to same timezone for comparison
		now := time.Now().In(wibLocation)

		// DB stores UTC, convert to WIB
		endTimeLocal := session.EndTime.In(wibLocation)

		if now.Before(endTimeLocal) {
			continue
		}

		// get bid data from redis
		bid, err := s.redisRepo.GetBidByKey(key)
		if err != nil {
			continue
		}

		// save final bid to DB
		err = s.bidRepo.SaveFinalBid(&entity.Bid{
			ItemID: itemID,
			UserID: bid.UserID,
			Amount: bid.Amount,
		})
		if err != nil {
			s.logger.Error("failed to save final bid", "sessionID", parsedSessionID, "itemID", itemID, "error", err)
			continue
		}

		// set status to finished
		item, err := s.itemRepo.GetByID(itemID)
		if err != nil {
			s.logger.Warn("failed to get item for status update", "itemID", itemID, "error", err)
		} else {
			item.Status = "finished"
			if err := s.itemRepo.Update(item); err != nil {
				s.logger.Error("failed to update item status", "itemID", itemID, "error", err)
			}
		}

		// delete redis key after saved to table Bid
		if err := s.redisRepo.DeleteKey(key); err != nil {
			s.logger.Warn("failed to delete Redis key", "key", key, "error", err)
		}

		s.logger.Info("final bid saved",
			"sessionID", parsedSessionID,
			"itemID", itemID,
			"amount", bid.Amount,
			"winner", bid.UserID,
		)
		totalSavedAuctionItemToDB++
	}

	if totalSavedAuctionItemToDB > 0 {
		s.logger.Info("expired sessions processed", "totalSaved", totalSavedAuctionItemToDB)
	}

	// also check for ongoing items with expired sessions with no bids
	if err := s.CloseExpiredItemsWithoutBids(); err != nil {
		s.logger.Error("failed to close expired items without bids", "error", err)
	}

	return nil
}

// CloseExpiredItemsWithoutBids to get the redis key with no bids
func (s *bidService) CloseExpiredItemsWithoutBids() error {
	items, err := s.itemRepo.GetAll()
	if err != nil {
		return err
	}

	// Convert to same timezone for comparison
	now := time.Now().In(wibLocation)
	closedCount := 0

	for _, item := range items {
		// validate item is not ongoing
		if item.Status != "ongoing" {
			continue
		}

		// skip if no session
		if item.SessionID == nil {
			continue
		}

		// get session to check end time
		session, err := s.auctionSessionRepo.GetByID(*item.SessionID)
		if err != nil {
			continue
		}

		// DB stores UTC, convert to WIB
		endTimeLocal := session.EndTime.In(wibLocation)

		if now.After(endTimeLocal) {
			amount, _, err := s.redisRepo.GetHighestBid(*item.SessionID, item.ID)

			// if no bid get status back to 'scheduled'
			if err != nil || amount == 0 {
				item.Status = "scheduled"
				if err := s.itemRepo.Update(&item); err != nil {
					s.logger.Error("failed to revert item to scheduled", "itemID", item.ID, "error", err)
				} else {
					closedCount++
				}
			}
		}
	}

	if closedCount > 0 {
		s.logger.Info("reverted items to scheduled", "count", closedCount)
	}

	return nil
}

func (s *bidService) DeleteKeyValue() error {
	pattern := "active:auction:*:item:*"
	keys, err := s.redisRepo.ScanKeys(pattern)
	if err != nil {
		s.logger.Error("failed to scan keys for cleanup", "pattern", pattern, "error", err)
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	deletedCount := 0
	for _, key := range keys {
		if err := s.redisRepo.DeleteKey(key); err != nil {
			s.logger.Warn("failed to delete key during cleanup", "key", key, "error", err)
		} else {
			deletedCount++
		}
	}

	s.logger.Info("midnight cleanup completed", "deleted", deletedCount)
	return nil
}
