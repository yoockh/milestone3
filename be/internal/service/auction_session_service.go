package service

import (
	"errors"
	"log/slog"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
	"time"
)

type sessionService struct {
	repo      repository.AuctionSessionRepository
	// redisRepo repository.SessionRedisRepository
	logger    *slog.Logger
}

type AuctionSessionService interface {
	Create(session *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error)
	GetByID(id int64) (dto.AuctionSessionDTO, error)
	GetAll() ([]dto.AuctionSessionDTO, error)
	Update(id int64, session *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error)
	Delete(id int64) error
}

func NewAuctionSessionService(r repository.AuctionSessionRepository, rr repository.SessionRedisRepository, logger *slog.Logger) AuctionSessionService {
	return &sessionService{repo: r, redisRepo: rr, logger: logger}
}

func (s *sessionService) Create(d *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error) {
	if d.Name == "" || d.EndTime.Before(d.StartTime) {
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	session := entity.AuctionSession{
		Name:      d.Name,
		StartTime: d.StartTime,
		EndTime:   d.EndTime,
	}

	err := s.repo.Create(&session)
	if err != nil {
		s.logger.Error("Failed to create auction session", "error", err)
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	now := time.Now()
	if session.StartTime.Before(now) && session.EndTime.After(now) {
		if err := s.redisRepo.SetActiveSession(session); err != nil {
			s.logger.Error("Failed to cache active session in Redis", "error", err)
		}
	}

	return dto.AuctionSessionResponse(session), nil
}

func (s *sessionService) GetByID(id int64) (dto.AuctionSessionDTO, error) {
	session, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction session by ID", "error", err)
		return dto.AuctionSessionDTO{}, ErrAuctionNotFound
	}

	return dto.AuctionSessionResponse(*session), nil
}

func (s *sessionService) GetAll() ([]dto.AuctionSessionDTO, error) {
	sessions, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all auction sessions", "error", err)
		return nil, ErrAuctionNotFound
	}

	var sessionDTOs []dto.AuctionSessionDTO
	for _, session := range sessions {
		sessionDTOs = append(sessionDTOs, dto.AuctionSessionResponse(*session))
	}

	return sessionDTOs, nil
}

func (s *sessionService) Update(id int64, d *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error) {
	session, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction session by ID for update", "error", err)
		return dto.AuctionSessionDTO{}, ErrAuctionNotFoundID
	}

	now := time.Now()
	if session.StartTime.Before(now) && session.EndTime.After(now) {
		return dto.AuctionSessionDTO{}, errors.New("cannot update an active session")
	}
	if session.EndTime.Before(now) {
		return dto.AuctionSessionDTO{}, errors.New("session already ended")
	}

	if d.Name != "" {
		session.Name = d.Name
	}
	if !d.StartTime.IsZero() {
		session.StartTime = d.StartTime
	}
	if !d.EndTime.IsZero() {
		session.EndTime = d.EndTime
	}

	if session.EndTime.Before(session.StartTime) {
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	// assign to DB
	err = s.repo.Update(session)
	if err != nil {
		s.logger.Error("Failed to update auction session", "error", err)
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	// update Redis if session is active
	if session.StartTime.Before(now) && session.EndTime.After(now) {
		if err := s.redisRepo.SetActiveSession(*session); err != nil {
			s.logger.Error("Failed to update active session in Redis", "error", err)
		}
	} else {
		if err := s.redisRepo.DeleteSession(session.ID); err != nil {
			s.logger.Error("Failed to delete session from Redis", "error", err)
		}
	}

	return dto.AuctionSessionResponse(*session), nil
}

func (s *sessionService) Delete(id int64) error {
	session, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction session by ID for delete", "error", err)
		return ErrAuctionNotFoundID
	}

	now := time.Now()

	if session.StartTime.Before(now) && session.EndTime.After(now) {
		// TODO: error contract
		return errors.New("cannot delete an active session")
	}

	if session.EndTime.Before(now) {
		// TODO: error contract
		return errors.New("session already ended")
	}

	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete auction session", "error", err)
		return ErrInvalidAuction
	}

	if err := s.redisRepo.DeleteSession(session.ID); err != nil {
		s.logger.Error("Failed to delete session from Redis", "error", err)
	}

	return nil
}
