package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type BidRepository interface {
	SaveFinalBid(bid *entity.Bid) error
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) SaveFinalBid(bid *entity.Bid) error {
	return r.db.Create(bid).Error
}
