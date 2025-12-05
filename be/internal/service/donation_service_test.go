package service

import (
	"errors"
	"testing"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDonationService_CreateDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		req     dto.DonationDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful donation creation",
			req: dto.DonationDTO{
				Title:       "Test Donation",
				Description: "Test description",
				Category:    "Electronics",
			},
			setup: func() {
				mockRepo.EXPECT().CreateDonation(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository create error",
			req: dto.DonationDTO{
				Title:       "Test Donation",
				Description: "Test description",
				Category:    "Electronics",
			},
			setup: func() {
				mockRepo.EXPECT().CreateDonation(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := donationService.CreateDonation(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDonationService_GetAllDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		userID  uint
		isAdmin bool
		setup   func()
		wantErr bool
	}{
		{
			name:    "admin get all donations",
			userID:  1,
			isAdmin: true,
			setup: func() {
				donations := []entity.Donation{
					{ID: 1, Title: "Donation 1", UserID: 1},
					{ID: 2, Title: "Donation 2", UserID: 2},
				}
				mockRepo.EXPECT().GetAllDonations(1, 10).Return(donations, int64(2), nil)
			},
			wantErr: false,
		},
		{
			name:    "user get own donations",
			userID:  1,
			isAdmin: false,
			setup: func() {
				donations := []entity.Donation{
					{ID: 1, Title: "Donation 1", UserID: 1},
				}
				mockRepo.EXPECT().GetDonationsByUserID(uint(1), 1, 10).Return(donations, int64(1), nil)
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			userID:  1,
			isAdmin: true,
			setup: func() {
				mockRepo.EXPECT().GetAllDonations(1, 10).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, total, err := donationService.GetAllDonations(tt.userID, tt.isAdmin, 1, 10)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, total, int64(0))
			}
		})
	}
}

func TestDonationService_GetDonationByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		id      uint
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get donation by id",
			id:   1,
			setup: func() {
				donation := entity.Donation{ID: 1, Title: "Test Donation", UserID: 1}
				mockRepo.EXPECT().GetDonationByID(uint(1)).Return(donation, nil)
			},
			wantErr: false,
		},
		{
			name: "donation not found",
			id:   999,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(999)).Return(entity.Donation{}, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := donationService.GetDonationByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), result.ID)
			}
		})
	}
}

func TestDonationService_CanManageDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name     string
		userID   uint
		ownerID  uint
		isAdmin  bool
		expected bool
	}{
		{
			name:     "admin can manage any donation",
			userID:   1,
			ownerID:  2,
			isAdmin:  true,
			expected: true,
		},
		{
			name:     "user can manage own donation",
			userID:   1,
			ownerID:  1,
			isAdmin:  false,
			expected: true,
		},
		{
			name:     "user cannot manage other's donation",
			userID:   1,
			ownerID:  2,
			isAdmin:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := donationService.CanManageDonations(tt.userID, tt.ownerID, tt.isAdmin)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDonationService_UpdateDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		req     dto.DonationDTO
		userID  uint
		isAdmin bool
		setup   func()
		wantErr bool
	}{
		{
			name: "successful update by owner",
			req: dto.DonationDTO{
				ID:    1,
				Title: "Updated",
			},
			userID:  1,
			isAdmin: false,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(1)).Return(entity.Donation{ID: 1, UserID: 1}, nil)
				mockRepo.EXPECT().UpdateDonation(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "forbidden - not owner",
			req: dto.DonationDTO{
				ID: 1,
			},
			userID:  2,
			isAdmin: false,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(1)).Return(entity.Donation{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := donationService.UpdateDonation(tt.req, tt.userID, tt.isAdmin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDonationService_DeleteDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		id      uint
		userID  uint
		isAdmin bool
		setup   func()
		wantErr bool
	}{
		{
			name:    "successful delete by owner",
			id:      1,
			userID:  1,
			isAdmin: false,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(1)).Return(entity.Donation{ID: 1, UserID: 1}, nil)
				mockRepo.EXPECT().DeleteDonation(uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "donation not found",
			id:      999,
			userID:  1,
			isAdmin: false,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(999)).Return(entity.Donation{}, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := donationService.DeleteDonation(tt.id, tt.userID, tt.isAdmin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDonationService_PatchDonation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDonationRepo(ctrl)
	mockStorage := mocks.NewMockGCPStorageRepo(ctrl)
	donationService := NewDonationService(mockRepo, mockStorage)

	tests := []struct {
		name    string
		req     dto.DonationDTO
		userID  uint
		isAdmin bool
		setup   func()
		wantErr bool
	}{
		{
			name: "successful patch by admin",
			req: dto.DonationDTO{
				ID:     1,
				Status: "verified_for_donation",
			},
			userID:  1,
			isAdmin: true,
			setup: func() {
				mockRepo.EXPECT().GetDonationByID(uint(1)).Return(entity.Donation{ID: 1, UserID: 2}, nil)
				mockRepo.EXPECT().PatchDonation(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := donationService.PatchDonation(tt.req, tt.userID, tt.isAdmin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
