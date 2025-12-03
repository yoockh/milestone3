package service

import (
	"context"
	"errors"
	"io"
	"time"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"

	"gorm.io/gorm"
)

type DonationService interface {
	CreateDonation(donationDTO dto.DonationDTO) error
	GetAllDonations(userID uint, isAdmin bool) ([]dto.DonationDTO, error)
	GetDonationByID(id uint) (dto.DonationDTO, error)
	UpdateDonation(donationDTO dto.DonationDTO, userID uint, isAdmin bool) error
	DeleteDonation(id uint, userID uint, isAdmin bool) error
	PatchDonation(donationDTO dto.DonationDTO, userID uint, isAdmin bool) error
	CanManageDonations(userID uint, ownerID uint, isAdmin bool) bool
}

type donationService struct {
	repo         repository.DonationRepo
	privateStore repository.GCPStorageRepo
}

func NewDonationService(repo repository.DonationRepo, privateStore repository.GCPStorageRepo) DonationService {
	return &donationService{
		repo:         repo,
		privateStore: privateStore,
	}
}

func (s *donationService) CreateDonation(donationDTO dto.DonationDTO) error {
	donation, err := dto.DonationRequest(donationDTO)
	if err != nil {
		return err
	}
	return s.repo.CreateDonation(donation)
}

func (s *donationService) GetAllDonations(userID uint, isAdmin bool) ([]dto.DonationDTO, error) {
	if isAdmin {
		donations, err := s.repo.GetAllDonations()
		if err != nil {
			return nil, err
		}
		return dto.DonationResponses(donations), nil
	}
	donations, err := s.repo.GetDonationsByUserID(userID)
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

func (s *donationService) PatchDonation(donationDTO dto.DonationDTO, userID uint, isAdmin bool) error {
	donation, err := dto.DonationRequest(donationDTO)
	if err != nil {
		return err
	}

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

	return s.repo.PatchDonation(donation)
}

func (s *donationService) CanManageDonations(userID uint, ownerID uint, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	return userID == ownerID
}

// ======================
//  METHODS FOR GCS
// ======================

func (s *donationService) UploadDonationImage(ctx context.Context, file io.Reader, fileName string) (string, error) {
	objectName, err := s.privateStore.UploadFile(ctx, file, fileName)
	if err != nil {
		return ErrImageNotFound.Error(), err
	}
	return objectName, nil
}

func (s *donationService) GetDonationImageURL(ctx context.Context, objectName string) (string, error) {
	url, err := s.privateStore.GenerateSignedURL(ctx, objectName, 10*time.Minute)
	if err != nil {
		return ErrSignedURLFailed.Error(), err
	}
	return url, nil
}
