package service

import (
	"errors"
	"log/slog"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type bidService struct {
	redisRepo repository.BidRedisRepository
	bidRepo   repository.BidRepository
	itemRepo  repository.AuctionItemRepository
	logger    *slog.Logger
}

type BidService interface {
	PlaceBid(sessionID, itemID, userID int64, amount float64, sessionEndTime time.Time) error
	GetHighestBid(sessionID, itemID int64) (float64, int64, error)
	SyncHighestBid(sessionID, itemID int64) error

	SaveExpiredSessions(sessionID int64) error
}

func NewBidService(r repository.BidRedisRepository, b repository.BidRepository, itemRepo repository.AuctionItemRepository, logger *slog.Logger) BidService {
	return &bidService{redisRepo: r, bidRepo: b, itemRepo: itemRepo, logger: logger}
}

func (s *bidService) PlaceBid(sessionID, itemID, userID int64, amount float64, sessionEndTime time.Time) error {
	if amount <= 0 {
		return ErrInvalidBidding
	}

	item, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		s.logger.Error("item not found", "itemID", itemID, "error", err)
		return ErrAuctionNotFound
	}

	s.logger.Info("item found", "itemID", itemID, "status", item.Status, "sessionID", item.SessionID)

	if item.Status != "ongoing" {
		s.logger.Warn("item not ongoing", "itemID", itemID, "status", item.Status)
		return ErrInvalidAuction
	}

	highest, bidder, err := s.redisRepo.GetHighestBid(sessionID, itemID)
	if err != nil && !errors.Is(err, redis.Nil) {
		s.logger.Error("failed to get highest bid", "error", err)
		return err
	}

	if amount <= highest {
		return ErrBidTooLow
	}

	err = s.redisRepo.SetHighestBid(sessionID, itemID, amount, userID, sessionEndTime)
	if err != nil {
		s.logger.Error("failed to set highest bid", "error", err)
		return err
	}

	s.logger.Info("new highest bid",
		"itemID", itemID,
		"userID", userID,
		"oldBid", highest,
		"oldBidder", bidder,
		"newBid", amount,
	)

	return nil
}

func (s *bidService) GetHighestBid(sessionID, itemID int64) (float64, int64, error) {
	_, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		return 0, 0, ErrAuctionNotFound
	}
	return s.redisRepo.GetHighestBid(sessionID, itemID)
}

func (s *bidService) SyncHighestBid(sessionID, itemID int64) error {
	highest, bidder, err := s.redisRepo.GetHighestBid(sessionID, itemID)
	if err != nil {
		s.logger.Error("failed to get highest bid", "error", err)
		return err
	}

	if highest == 0 || bidder == 0 {
		return nil
	}

	finalBid := &entity.Bid{
		SessionID: sessionID,
		ItemID:    itemID,
		Amount:    highest,
		UserID:    bidder,
	}

	err = s.bidRepo.SaveFinalBid(finalBid)
	if err != nil {
		s.logger.Error("failed to save final bid to postgres", "error", err)
		return err
	}

	s.logger.Info("highest bid synced to postgres",
		"sessionID", sessionID,
		"itemID", itemID,
		"bid", highest,
		"bidder", bidder,
	)

	return nil
}

func parseKey(key string) (sessionID, itemID int64, err error) {
	parts := strings.Split(key, ":")
	// format: auction_session:{sessionID}:item:{itemID}
	sID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	iID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return sID, iID, nil
}

func (s *bidService) SaveExpiredSessions(sessionID int64) error {
	keys, err := s.redisRepo.ScanKeys("auction_session:*:item:*")
	if err != nil {
		return err
	}

	for _, key := range keys {
		// get key: auction_session:10:item:5
		sessionID, itemID, err := parseKey(key)
		if err != nil {
			continue
		}

		// result from bid redis
		bid, err := s.redisRepo.GetBidByKey(key)
		if err != nil {
			continue
		}

		// get end time
		endTime, err := s.redisRepo.GetSessionEndTime(sessionID)
		if err != nil {
			continue
		}

		if time.Now().Before(endTime) {
			continue // belum expired
		}

		// assign to DB
		err = s.bidRepo.SaveFinalBid(&entity.Bid{
			SessionID: sessionID,
			ItemID:    itemID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
		})
		if err != nil {
			continue
		}

		_ = s.redisRepo.DeleteKey(key)
	}

	return nil
}
