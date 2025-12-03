package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type AuctionItemRepository interface {
	Create(item *entity.AuctionItem) error
	GetAll() ([]entity.AuctionItem, error)
	GetByID(id int64) (*entity.AuctionItem, error)
	ReadBySession(sessionID int64) ([]entity.AuctionItem, error)
	Update(item *entity.AuctionItem) error
	Delete(id int64) error
}

type auctionItemRepository struct {
	db *gorm.DB
}

func NewAuctionItemRepository(db *gorm.DB) AuctionItemRepository {
	return &auctionItemRepository{db: db}
}

func (r *auctionItemRepository) Create(item *entity.AuctionItem) error {
	return r.db.Create(item).Error
}

func (r *auctionItemRepository) GetAll() ([]entity.AuctionItem, error) {
	var items []entity.AuctionItem
	err := r.db.Preload("Session").Preload("Photos").Find(&items).Error
	return items, err
}

func (r *auctionItemRepository) GetByID(id int64) (*entity.AuctionItem, error) {
	var item entity.AuctionItem
	err := r.db.Preload("Session").First(&item, id).Error
	return &item, err
}

func (r *auctionItemRepository) ReadBySession(sessionID int64) ([]entity.AuctionItem, error) {
	var items []entity.AuctionItem
	err := r.db.Preload("Session").Preload("Photos").Where("session_id = ?", sessionID).Find(&items).Error
	return items, err
}

func (r *auctionItemRepository) Update(item *entity.AuctionItem) error {
	return r.db.Save(item).Error
}

func (r *auctionItemRepository) Delete(id int64) error {
	return r.db.Delete(&entity.AuctionItem{}, id).Error
}
