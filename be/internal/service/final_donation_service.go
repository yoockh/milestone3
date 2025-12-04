package service

import (
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
)

type FinalDonationService interface {
	GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error)
	GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error)
}

type finalDonationService struct {
	finalDonationRepo repository.FinalDonationRepository
}

func NewFinalDonationService(finalDonationRepo repository.FinalDonationRepository) FinalDonationService {
	return &finalDonationService{finalDonationRepo: finalDonationRepo}
}

func (s *finalDonationService) GetAllFinalDonations(page, limit int) ([]entity.FinalDonation, int64, error) {
	return s.finalDonationRepo.GetAllFinalDonations(page, limit)
}

func (s *finalDonationService) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	return s.finalDonationRepo.GetAllFinalDonationsByUserID(userID)
}
