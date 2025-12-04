package repository

import (
	"milestone3/be/internal/entity"

	"gorm.io/gorm"
)

type FinalDonationRepository interface {
	GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error)
	GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error)
}

type finalDonationRepository struct {
	db *gorm.DB
}

func NewFinalDonationRepository(db *gorm.DB) FinalDonationRepository {
	return &finalDonationRepository{db: db}
}

// Return final_donations where the related donation has status = entity.StatusVerifiedForDonation
func (r *finalDonationRepository) GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error) {
	var finalDonations []entity.FinalDonation
	var total int64

	// Count total records
	if err := r.db.Model(&entity.FinalDonation{}).
		Joins("JOIN donations d ON d.id = final_donations.donation_id").
		Where("d.status = ?", entity.StatusVerifiedForDonation).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * limit
	err := r.db.
		Joins("JOIN donations d ON d.id = final_donations.donation_id").
		Where("d.status = ?", entity.StatusVerifiedForDonation).
		Preload("Donation").
		Offset(offset).Limit(limit).
		Order("final_donations.created_at DESC").
		Find(&finalDonations).Error
	return finalDonations, total, err
}

// Return final_donations for a user by joining donations and filtering by donation.user_id
func (r *finalDonationRepository) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	var finalDonations []entity.FinalDonation
	err := r.db.
		Joins("JOIN donations d ON d.id = final_donations.donation_id").
		Where("d.user_id = ? AND d.status = ?", userID, entity.StatusVerifiedForDonation).
		Preload("Donation").
		Find(&finalDonations).Error
	return finalDonations, err
}
