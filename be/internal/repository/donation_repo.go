package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type DonationRepo interface {
	CreateDonation(donation entity.Donation) error
	GetDonationByID(id uint) (entity.Donation, error)
	UpdateDonation(donation entity.Donation) error
	DeleteDonation(id uint) error

	// Admin-only or filtered queries with pagination
	GetAllDonations(page, limit int) ([]entity.Donation, int64, error)
	GetDonationsByUserID(userID uint, page, limit int) ([]entity.Donation, int64, error)

	PatchDonation(donation entity.Donation) error
}

type donationRepo struct {
	db *gorm.DB
}

func NewDonationRepo(db *gorm.DB) DonationRepo {
	return &donationRepo{db: db}
}

func (r *donationRepo) CreateDonation(donation entity.Donation) error {
	return r.db.Create(&donation).Error
}

func (r *donationRepo) GetAllDonations(page, limit int) ([]entity.Donation, int64, error) {
	var donations []entity.Donation
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Donation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with preload
	offset := (page - 1) * limit
	err := r.db.Preload("Photos").Offset(offset).Limit(limit).Order("created_at DESC").Find(&donations).Error
	return donations, total, err
}

func (r *donationRepo) GetDonationsByUserID(userID uint, page, limit int) ([]entity.Donation, int64, error) {
	var donations []entity.Donation
	var total int64

	// Count total records for user
	if err := r.db.Model(&entity.Donation{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with preload
	offset := (page - 1) * limit
	err := r.db.Preload("Photos").Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&donations).Error
	return donations, total, err
}

func (r *donationRepo) GetDonationByID(id uint) (entity.Donation, error) {
	var donation entity.Donation
	err := r.db.Preload("Photos").First(&donation, id).Error
	return donation, err
}

func (r *donationRepo) UpdateDonation(donation entity.Donation) error {
	return r.db.Save(&donation).Error
}

func (r *donationRepo) DeleteDonation(id uint) error {
	return r.db.Delete(&entity.Donation{}, id).Error
}

func (r *donationRepo) PatchDonation(donation entity.Donation) error {
	return r.db.Model(&entity.Donation{}).Where("id = ?", donation.ID).Updates(donation).Error
}
