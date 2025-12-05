package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuctionSessionService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionSessionRepository(ctrl)
	mockRedis := mocks.NewMockSessionRedisRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	sessionService := NewAuctionSessionService(mockRepo, mockRedis, logger)

	now := time.Now()
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		req     dto.AuctionSessionDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful session creation",
			req: dto.AuctionSessionDTO{
				Name:      "Test Session",
				StartTime: now,
				EndTime:   future,
			},
			setup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
				mockRedis.EXPECT().SetActiveSession(gomock.Any()).Return(nil).AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "invalid session - empty name",
			req: dto.AuctionSessionDTO{
				Name:      "",
				StartTime: now,
				EndTime:   future,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "invalid session - end before start",
			req: dto.AuctionSessionDTO{
				Name:      "Test Session",
				StartTime: future,
				EndTime:   now,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "repository create error",
			req: dto.AuctionSessionDTO{
				Name:      "Test Session",
				StartTime: now,
				EndTime:   future,
			},
			setup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			result, err := sessionService.Create(&tt.req)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Test Session", result.Name)
			}
		})
	}
}

func TestAuctionSessionService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionSessionRepository(ctrl)
	mockRedis := mocks.NewMockSessionRedisRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	sessionService := NewAuctionSessionService(mockRepo, mockRedis, logger)

	tests := []struct {
		name    string
		id      int64
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get by id",
			id:   1,
			setup: func() {
				session := &entity.AuctionSession{
					ID:   1,
					Name: "Test Session",
				}
				mockRepo.EXPECT().GetByID(int64(1)).Return(session, nil)
			},
			wantErr: false,
		},
		{
			name: "session not found",
			id:   999,
			setup: func() {
				mockRepo.EXPECT().GetByID(int64(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			result, err := sessionService.GetByID(tt.id)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), int64(result.ID))
			}
		})
	}
}

func TestAuctionSessionService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionSessionRepository(ctrl)
	mockRedis := mocks.NewMockSessionRedisRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	sessionService := NewAuctionSessionService(mockRepo, mockRedis, logger)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get all sessions",
			setup: func() {
				sessions := []*entity.AuctionSession{
					{ID: 1, Name: "Session 1"},
					{ID: 2, Name: "Session 2"},
				}
				mockRepo.EXPECT().GetAll().Return(sessions, nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func() {
				mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			
			result, err := sessionService.GetAll()
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
			}
		})
	}
}

func TestAuctionSessionService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionSessionRepository(ctrl)
	mockRedis := mocks.NewMockSessionRedisRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sessionService := NewAuctionSessionService(mockRepo, mockRedis, logger)

	tests := []struct {
		name    string
		id      int64
		req     *dto.AuctionSessionDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful update",
			id:   1,
			req:  &dto.AuctionSessionDTO{Name: "Updated"},
			setup: func() {
				futureTime := time.Now().Add(24 * time.Hour)
				session := &entity.AuctionSession{ID: 1, Name: "Old", StartTime: futureTime, EndTime: futureTime.Add(2 * time.Hour)}
				mockRepo.EXPECT().GetByID(int64(1)).Return(session, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
				mockRedis.EXPECT().DeleteSession(int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update active session - error",
			id:   1,
			req:  &dto.AuctionSessionDTO{Name: "Updated"},
			setup: func() {
				pastTime := time.Now().Add(-1 * time.Hour)
				futureTime := time.Now().Add(1 * time.Hour)
				session := &entity.AuctionSession{ID: 1, StartTime: pastTime, EndTime: futureTime}
				mockRepo.EXPECT().GetByID(int64(1)).Return(session, nil)
			},
			wantErr: true,
		},
		{
			name: "session not found",
			id:   999,
			req:  &dto.AuctionSessionDTO{},
			setup: func() {
				mockRepo.EXPECT().GetByID(int64(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := sessionService.Update(tt.id, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuctionSessionService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionSessionRepository(ctrl)
	mockRedis := mocks.NewMockSessionRedisRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sessionService := NewAuctionSessionService(mockRepo, mockRedis, logger)

	tests := []struct {
		name    string
		id      int64
		setup   func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   1,
			setup: func() {
				futureTime := time.Now().Add(24 * time.Hour)
				session := &entity.AuctionSession{ID: 1, StartTime: futureTime, EndTime: futureTime.Add(2 * time.Hour)}
				mockRepo.EXPECT().GetByID(int64(1)).Return(session, nil)
				mockRepo.EXPECT().Delete(int64(1)).Return(nil)
				mockRedis.EXPECT().DeleteSession(int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete active session - error",
			id:   1,
			setup: func() {
				pastTime := time.Now().Add(-1 * time.Hour)
				futureTime := time.Now().Add(1 * time.Hour)
				session := &entity.AuctionSession{ID: 1, StartTime: pastTime, EndTime: futureTime}
				mockRepo.EXPECT().GetByID(int64(1)).Return(session, nil)
			},
			wantErr: true,
		},
		{
			name: "session not found",
			id:   999,
			setup: func() {
				mockRepo.EXPECT().GetByID(int64(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := sessionService.Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
