package service

import (
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
)

type FinalDonationService interface {
	GetAllFinalDonations() ([]entity.FinalDonation, error)
	GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error)
}

type finalDonationService struct {
	finalDonationRepo repository.FinalDonationRepository
}

func NewFinalDonationService(finalDonationRepo repository.FinalDonationRepository) FinalDonationService {
	return &finalDonationService{finalDonationRepo: finalDonationRepo}
}

func (s *finalDonationService) GetAllFinalDonations() ([]entity.FinalDonation, error) {
	return s.finalDonationRepo.GetAllFinalDonations()
}

func (s *finalDonationService) GetAllFinalDonationsByUserID(userID int) ([]entity.FinalDonation, error) {
	return s.finalDonationRepo.GetAllFinalDonationsByUserID(userID)
}
