package service

import (
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
)

type FinalDonationService interface {
	GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error)
	GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error)
	UpdateNotes(donationID uint, userID uint, notes string) error
}

type finalDonationService struct {
	finalDonationRepo repository.FinalDonationRepository
	donationRepo      repository.DonationRepo
}

func NewFinalDonationService(finalDonationRepo repository.FinalDonationRepository, donationRepo repository.DonationRepo) FinalDonationService {
	return &finalDonationService{
		finalDonationRepo: finalDonationRepo,
		donationRepo:      donationRepo,
	}
}

func (s *finalDonationService) GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error) {
	return s.finalDonationRepo.GetAllFinalDonations(page, limit)
}

func (s *finalDonationService) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	return s.finalDonationRepo.GetAllFinalDonationsByUserID(userID)
}

func (s *finalDonationService) UpdateNotes(donationID uint, userID uint, notes string) error {
	// Get donation to check ownership and status
	donation, err := s.donationRepo.GetDonationByID(donationID)
	if err != nil {
		return ErrDonationNotFound
	}

	// Check ownership first
	if donation.UserID != userID {
		return ErrForbidden
	}

	// Check if status is verified_for_donation
	if donation.Status != entity.StatusVerifiedForDonation {
		return ErrDonationNotVerified
	}

	return s.finalDonationRepo.UpdateNotes(donationID, notes)
}
