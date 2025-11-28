package service

import (
	"errors"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"

	"gorm.io/gorm"
)

type DonationService interface {
	CreateDonation(donationDTO dto.DonationDTO) error
	GetAllDonations() ([]dto.DonationDTO, error)
	GetDonationByID(id uint) (dto.DonationDTO, error)
	// Update/Delete now require caller identity: userID and isAdmin flag
	UpdateDonation(donationDTO dto.DonationDTO, userID uint, isAdmin bool) error
	DeleteDonation(id uint, userID uint, isAdmin bool) error

	// helper to check permission
	CanManageDonations(userID uint, ownerID uint, isAdmin bool) bool
}

type donationService struct {
	repo repository.DonationRepo
}

func NewDonationService(repo repository.DonationRepo) DonationService {
	return &donationService{repo: repo}
}

func (s *donationService) CreateDonation(donationDTO dto.DonationDTO) error {
	donation, err := dto.DonationRequest(donationDTO)
	if err != nil {
		return err
	}
	return s.repo.CreateDonation(donation)
}

func (s *donationService) GetAllDonations() ([]dto.DonationDTO, error) {
	donations, err := s.repo.GetAllDonations()
	if err != nil {
		return nil, err
	}
	return dto.DonationResponses(donations), nil
}

func (s *donationService) GetDonationByID(id uint) (dto.DonationDTO, error) {
	donation, err := s.repo.GetDonationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DonationDTO{}, ErrDonationNotFound
		}
		return dto.DonationDTO{}, err
	}
	return dto.DonationResponse(donation), nil
}

func (s *donationService) UpdateDonation(donationDTO dto.DonationDTO, userID uint, isAdmin bool) error {
	donation, err := dto.DonationRequest(donationDTO)
	if err != nil {
		return err
	}

	// verify existing donation and owner
	existing, err := s.repo.GetDonationByID(donation.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDonationNotFound
		}
		return err
	}

	if !s.CanManageDonations(userID, existing.UserID, isAdmin) {
		return ErrForbidden
	}

	return s.repo.UpdateDonation(donation)
}

func (s *donationService) DeleteDonation(id uint, userID uint, isAdmin bool) error {
	// fetch to check owner
	donation, err := s.repo.GetDonationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDonationNotFound
		}
		return err
	}

	if !s.CanManageDonations(userID, donation.UserID, isAdmin) {
		return ErrForbidden
	}

	return s.repo.DeleteDonation(id)
}

func (s *donationService) CanManageDonations(userID uint, ownerID uint, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	return userID == ownerID
}
