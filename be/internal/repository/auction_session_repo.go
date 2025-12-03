package repository

import (
	"milestone3/be/internal/entity"
	"time"

	"gorm.io/gorm"
)

type AuctionSessionRepository interface {
	Create(session *entity.AuctionSession) error
	GetByID(id int64) (*entity.AuctionSession, error)
	GetAll() ([]*entity.AuctionSession, error)
	GetActiveSessions() ([]*entity.AuctionSession, error)
	Update(session *entity.AuctionSession) error
	Delete(id int64) error
}

type auctionSessionRepository struct {
	db *gorm.DB
}

func NewAuctionSessionRepository(db *gorm.DB) AuctionSessionRepository {
	return &auctionSessionRepository{db: db}
}

func (r *auctionSessionRepository) Create(session *entity.AuctionSession) error {
	return r.db.Create(session).Error
}

func (r *auctionSessionRepository) GetByID(id int64) (*entity.AuctionSession, error) {
	var session entity.AuctionSession
	err := r.db.First(&session, id).Error
	return &session, err
}

func (r *auctionSessionRepository) GetAll() ([]*entity.AuctionSession, error) {
	var sessions []*entity.AuctionSession
	err := r.db.Find(&sessions).Error
	return sessions, err
}

func (r *auctionSessionRepository) GetActiveSessions() ([]*entity.AuctionSession, error) {
	var sessions []*entity.AuctionSession
	now := time.Now()
	err := r.db.Where("start_time <= ? AND end_time >= ?", now, now).Find(&sessions).Error
	return sessions, err
}

func (r *auctionSessionRepository) Update(session *entity.AuctionSession) error {
	return r.db.Save(session).Error
}

func (r *auctionSessionRepository) Delete(id int64) error {
	return r.db.Delete(&entity.AuctionSession{}, id).Error
}
