package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type BidRedisRepository interface {
	SetHighestBid(sessionID, itemID int64, amount float64, userID int64, sessionEndTime time.Time) error
	GetHighestBid(sessionID, itemID int64) (float64, int64, error)
	GetBidHistory(sessionID, itemID int64, limit int64) ([]BidEntry, error)
	GetSessionEndTime(sessionID int64) (time.Time, error)

	ScanKeys(pattern string) ([]string, error)
	GetBidByKey(key string) (BidEntry, error)
	DeleteKey(key string) error
}

type bidRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

type BidEntry struct {
	UserID int64
	ItemID int64
	Amount float64
}

func NewBidRedisRepository(client *redis.Client, ctx context.Context) BidRedisRepository {
	return &bidRedisRepository{client: client, ctx: ctx}
}

func (r *bidRedisRepository) SetHighestBid(sessionID, itemID int64, amount float64, userID int64, sessionEndTime time.Time) error {
	key := fmt.Sprintf("active:auction:%d:item:%d", sessionID, itemID)

	if err := r.client.HSet(r.ctx, key, map[string]interface{}{
		"highest_amount": amount,
		"highest_bidder": userID,
		"updated_at":     time.Now().Unix(),
		"end_time":       sessionEndTime.Unix(),
	}).Err(); err != nil {
		return err
	}

	ttl := time.Until(sessionEndTime)
	if ttl > 0 {
		if err := r.client.Expire(r.ctx, key, ttl).Err(); err != nil {
			return err
		}
	}

	historyKey := fmt.Sprintf("auction:%d:item:%d:history", sessionID, itemID)
	if err := r.client.ZAdd(r.ctx, historyKey, redis.Z{
		Score:  amount,
		Member: userID,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *bidRedisRepository) GetHighestBid(sessionID, itemID int64) (float64, int64, error) {
	key := fmt.Sprintf("active:auction:%d:item:%d", sessionID, itemID)

	data, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	amount, err := strconv.ParseFloat(data["highest_amount"], 64)
	if err != nil {
		amount = 0
	}
	bidder, err := strconv.ParseInt(data["highest_bidder"], 10, 64)
	if err != nil {
		bidder = 0
	}

	return amount, bidder, nil
}

func (r *bidRedisRepository) GetBidHistory(sessionID, itemID int64, limit int64) ([]BidEntry, error) {
	historyKey := fmt.Sprintf("auction:%d:item:%d:history", sessionID, itemID)

	results, err := r.client.ZRevRangeWithScores(r.ctx, historyKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	var history []BidEntry
	for _, z := range results {
		userID, _ := strconv.ParseInt(fmt.Sprintf("%v", z.Member), 10, 64)
		history = append(history, BidEntry{
			UserID: userID,
			ItemID: itemID,
			Amount: z.Score,
		})
	}

	return history, nil
}

func (r *bidRedisRepository) GetSessionEndTime(sessionID int64) (time.Time, error) {
	{
		key := "active_session:" + strconv.FormatInt(sessionID, 10)
		data, err := r.client.HGetAll(r.ctx, key).Result()
		if err != nil {
			return time.Time{}, err
		}
		endTime, err := strconv.ParseInt(data["EndTime"], 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(endTime, 0), nil
	}
}

func (r *bidRedisRepository) ScanKeys(pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		var k []string
		var err error
		k, cursor, err = r.client.Scan(r.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (r *bidRedisRepository) GetBidByKey(key string) (BidEntry, error) {
	data, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return BidEntry{}, err
	}
	amount, _ := strconv.ParseFloat(data["highest_amount"], 64)
	userID, _ := strconv.ParseInt(data["highest_bidder"], 10, 64)
	itemID, _ := strconv.ParseInt(data["auction_item"], 10, 64)
	return BidEntry{
		UserID: userID,
		ItemID: itemID,
		Amount: amount,
	}, nil
}

func (r *bidRedisRepository) DeleteKey(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
