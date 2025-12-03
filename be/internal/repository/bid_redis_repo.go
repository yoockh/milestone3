package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type BidRedisRepository interface {
	SetHighestBid(sessionID, itemID int64, amount float64, userID int64, sessionEndTime time.Time) error
	GetHighestBid(sessionID, itemID int64) (float64, int64, error)
	GetEndTime(key string) (time.Time, error)

	ScanKeys(pattern string) ([]string, error)
	GetBidByKey(key string) (BidEntry, error)
	DeleteKey(key string) error

	CheckDuplicateBid(userID, itemID int64, amount float64, ttl time.Duration) error
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

	if len(data) == 0 {
		return 0, 0, nil
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

func (r *bidRedisRepository) GetEndTime(key string) (time.Time, error) {
	data, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return time.Time{}, err
	}

	if len(data) == 0 {
		return time.Time{}, errors.New("key not found")
	}

	endTimeStr, exists := data["end_time"]
	if !exists {
		return time.Time{}, errors.New("end_time field not found")
	}

	endTimeUnix, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse end_time: %w", err)
	}

	return time.Unix(endTimeUnix, 0), nil
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

	var sessionID, itemID int64
	_, err = fmt.Sscanf(key, "active:auction:%d:item:%d", &sessionID, &itemID)
	if err != nil {
		return BidEntry{}, fmt.Errorf("invalid key format: %s", key)
	}

	return BidEntry{
		UserID: userID,
		ItemID: itemID,
		Amount: amount,
	}, nil
}

func (r *bidRedisRepository) DeleteKey(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *bidRedisRepository) CheckDuplicateBid(userID, itemID int64, amount float64, ttl time.Duration) error {
	key := fmt.Sprintf("bidder:%d:item:%d:amount:%.2f", userID, itemID, amount)
	result, err := r.client.SetNX(r.ctx, key, "exists", ttl).Result()
	if err != nil {
		return err
	}
	if !result {
		return errors.New("duplicate bid detected")
	}
	return nil
}
