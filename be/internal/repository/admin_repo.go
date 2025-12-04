package repository

import (
	"context"
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type AdminRepo struct {
	db *gorm.DB
	ctx context.Context
}

func NewAdminRepository(db *gorm.DB, ctx context.Context) *AdminRepo {
	return &AdminRepo{db: db, ctx: ctx}
}


//count total transaction
func (ar *AdminRepo) CountPayment() (count int64, err error) {
	var payment entity.Payment
	if err := ar.db.WithContext(ar.ctx).Model(&payment).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
 
// //count total donation
func (ar *AdminRepo) CountDonation() (count int64, err error) {
	var donation entity.Donation
	if err := ar.db.WithContext(ar.ctx).Model(&donation).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// count total auction
// work in progress (WIP)
func (ar *AdminRepo) CountAuction() (count int64, err error) {
	var auction entity.AuctionItem
	if err := ar.db.WithContext(ar.ctx).Model(&auction).Count(&count).Error; err != nil {
			return 0, err
	}

	return count, nil
}

// //count total article
func (ar *AdminRepo) CountArticle() (count int64, err error) {
	var article entity.Article
	if err := ar.db.WithContext(ar.ctx).Model(&article).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// for reporting endpoint //
// work in progress (WIP)