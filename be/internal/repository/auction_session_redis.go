package repository

import (
	"context"
	"milestone3/be/internal/entity"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRedisRepository interface {
	SetActiveSession(session entity.AuctionSession) error
	GetActiveSessions() ([]entity.AuctionSession, error)
	DeleteSession(id int64) error
}

type sessionRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewSessionRedisRepository(client *redis.Client, ctx context.Context) SessionRedisRepository {
	return &sessionRedisRepository{client: client, ctx: ctx}
}

func (r *sessionRedisRepository) SetActiveSession(session entity.AuctionSession) error {
	key := "active_session:" + strconv.FormatInt(session.ID, 10)
	return r.client.HSet(r.ctx, key, map[string]interface{}{
		"ID":        session.ID,
		"Name":      session.Name,
		"StartTime": session.StartTime.Unix(),
		"EndTime":   session.EndTime.Unix(),
	}).Err()
}

func (r *sessionRedisRepository) GetActiveSessions() ([]entity.AuctionSession, error) {
	var sessions []entity.AuctionSession

	iter := r.client.Scan(r.ctx, 0, "active_session:*", 0).Iterator()
	for iter.Next(r.ctx) {
		key := iter.Val()
		data, err := r.client.HGetAll(r.ctx, key).Result()
		if err != nil {
			return nil, err
		}

		id, _ := strconv.ParseInt(data["ID"], 10, 64)
		startTimeUnix, _ := strconv.ParseInt(data["StartTime"], 10, 64)
		endTimeUnix, _ := strconv.ParseInt(data["EndTime"], 10, 64)

		session := entity.AuctionSession{
			ID:        id,
			Name:      data["Name"],
			StartTime: time.Unix(startTimeUnix, 0),
			EndTime:   time.Unix(endTimeUnix, 0),
		}

		sessions = append(sessions, session)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRedisRepository) DeleteSession(id int64) error {
	key := "active_session:" + strconv.FormatInt(id, 10)
	return r.client.Del(r.ctx, key).Err()
}
