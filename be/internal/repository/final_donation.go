package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type FinalDonationRepository interface {
	GetAllFinalDonations() ([]entity.FinalDonation, error)
	GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error)
}

type finalDonationRepository struct {
	db *gorm.DB
}

func NewFinalDonationRepository(db *gorm.DB) FinalDonationRepository {
	return &finalDonationRepository{db: db}
}

func (r *finalDonationRepository) GetAllFinalDonations() ([]entity.FinalDonation, error) {
	var finalDonations []entity.FinalDonation
	if err := r.db.Find(&finalDonations).Error; err != nil {
		return nil, err
	}
	return finalDonations, nil
}

func (r *finalDonationRepository) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	var finalDonations []entity.FinalDonation
	if err := r.db.Where("user_id = ?", userID).Find(&finalDonations).Error; err != nil {
		return nil, err
	}
	return finalDonations, nil
}
