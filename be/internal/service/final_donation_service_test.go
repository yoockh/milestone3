package service

import (
	"errors"
	"testing"

	"milestone3/be/internal/entity"
	"milestone3/be/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFinalDonationService_GetAllFinalDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFinalDonationRepository(ctrl)
	mockDonationRepo := mocks.NewMockDonationRepo(ctrl)
	finalDonationService := NewFinalDonationService(mockRepo, mockDonationRepo)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get all final donations",
			setup: func() {
				donations := []entity.FinalDonation{
					{ID: 1, DonationID: 1, Notes: "Institution 1"},
					{ID: 2, DonationID: 2, Notes: "Institution 2"},
				}
				mockRepo.EXPECT().GetAllFinalDonations(1, 10).Return(donations, int64(2), nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func() {
				mockRepo.EXPECT().GetAllFinalDonations(1, 10).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, total, err := finalDonationService.GetAllFinalDonations(1, 10)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, int64(2), total)
				assert.Equal(t, "Institution 1", result[0].Notes)
			}
		})
	}
}
func TestFinalDonationService_GetAllFinalDonationsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFinalDonationRepository(ctrl)
	mockDonationRepo := mocks.NewMockDonationRepo(ctrl)
	finalDonationService := NewFinalDonationService(mockRepo, mockDonationRepo)

	tests := []struct {
		name    string
		userID  int
		setup   func()
		wantErr bool
	}{
		{
			name:   "successful get final donations by user id",
			userID: 1,
			setup: func() {
				donations := []entity.FinalDonation{
					{ID: 1, DonationID: 1, Notes: "Institution 1"},
				}
				mockRepo.EXPECT().GetAllFinalDonationsByUserID(1).Return(donations, nil)
			},
			wantErr: false,
		},
		{
			name:   "repository error",
			userID: 1,
			setup: func() {
				mockRepo.EXPECT().GetAllFinalDonationsByUserID(1).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "no donations found for user",
			userID: 999,
			setup: func() {
				mockRepo.EXPECT().GetAllFinalDonationsByUserID(999).Return([]entity.FinalDonation{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := finalDonationService.GetAllFinalDonationsByUserID(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		}) // Closing the t.Run function
	} // Closing the TestFinalDonationService_GetAllFinalDonationsByUserID function
}
