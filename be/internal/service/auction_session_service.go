package service

import (
	"log/slog"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/repository"
	"time"
)

type sessionService struct {
	repo   repository.AuctionSessionRepository
	logger *slog.Logger
}

type AuctionSessionService interface {
	Create(session *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error)
	GetByID(id int64) (dto.AuctionSessionDTO, error)
	GetAll() ([]dto.AuctionSessionDTO, error)
	Update(id int64, session *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error)
	Delete(id int64) error
}

func NewAuctionSessionService(r repository.AuctionSessionRepository, logger *slog.Logger) AuctionSessionService {
	return &sessionService{repo: r, logger: logger}
}

func (s *sessionService) Create(d *dto.AuctionSessionDTO) (dto.AuctionSessionDTO, error) {
	// validate title
	if d.Name == "" {
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	// set local time
	now := time.Now()
	startTime := d.StartTime
	endTime := d.EndTime

	// end time > start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return dto.AuctionSessionDTO{}, ErrInvalidDate
	}

	// start time cannot be in the past (-1 minute for buffer)
	if startTime.Before(now.Add(-1 * time.Minute)) {
		return dto.AuctionSessionDTO{}, ErrInvalidTime
	}

	// end time cannot be in the past
	if endTime.Before(now) {
		return dto.AuctionSessionDTO{}, ErrInvalidTime
	}

	// minimum session duration (at least 1 minute)
	minDuration := 1 * time.Minute
	duration := endTime.Sub(startTime)
	if duration < minDuration {
		return dto.AuctionSessionDTO{}, ErrInvalidDate
	}

	// maximum session duration (24 hours to prevent mistakes)
	maxDuration := 24 * time.Hour
	if duration > maxDuration {
		return dto.AuctionSessionDTO{}, ErrInvalidDate
	}

	// Convert to UTC before saving to database
	// PostgreSQL TIMESTAMP (without timezone) stores as-is, so we convert to UTC first
	// This ensures consistency: client sends WIB (+07:00) -> we convert to UTC -> store as UTC
	session := entity.AuctionSession{
		Name:      d.Name,
		StartTime: d.StartTime.UTC(),
		EndTime:   d.EndTime.UTC(),
	}

	err := s.repo.Create(&session)
	if err != nil {
		s.logger.Error("Failed to create auction session", "error", err)
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	s.logger.Info("Auction session created", "sessionID", session.ID, "name", session.Name)
	return dto.AuctionSessionResponse(session), nil
}

func (s *sessionService) GetByID(id int64) (dto.AuctionSessionDTO, error) {
	session, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction session by ID", "error", err)
		return dto.AuctionSessionDTO{}, ErrSessionNotFoundID
	}

	return dto.AuctionSessionResponse(*session), nil
}

func (s *sessionService) GetAll() ([]dto.AuctionSessionDTO, error) {
	sessions, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all auction sessions", "error", err)
		return nil, ErrSessionNotFoundID
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
		return dto.AuctionSessionDTO{}, ErrSessionNotFoundID
	}

	// Convert DB times (stored as UTC) to local for comparison
	now := time.Now().In(wibLocation)
	sessionStart := session.StartTime.In(wibLocation)
	sessionEnd := session.EndTime.In(wibLocation)

	// cannot update active (ongoing) session
	if sessionStart.Before(now) && sessionEnd.After(now) {
		return dto.AuctionSessionDTO{}, ErrActiveSession
	}

	// cannot update expired session
	if sessionEnd.Before(now) {
		return dto.AuctionSessionDTO{}, ErrExpiredSession
	}

	// update fields
	if d.Name != "" {
		session.Name = d.Name
	}
	if !d.StartTime.IsZero() {
		// Convert to UTC before saving
		session.StartTime = d.StartTime.UTC()
	}
	if !d.EndTime.IsZero() {
		// Convert to UTC before saving
		session.EndTime = d.EndTime.UTC()
	}

	newStart := session.StartTime
	newEnd := session.EndTime

	// end time > start time
	if newEnd.Before(newStart) || newEnd.Equal(newStart) {
		return dto.AuctionSessionDTO{}, ErrInvalidDate
	}

	// start time cannot be in the past (-1 minute for buffer)
	if newStart.Before(now.Add(-1 * time.Minute)) {
		return dto.AuctionSessionDTO{}, ErrInvalidTime
	}

	// end time cannot be in the past
	if newEnd.Before(now) {
		return dto.AuctionSessionDTO{}, ErrInvalidTime
	}

	// minimum session duration (at least 1 minute)
	minDuration := 1 * time.Minute
	if newEnd.Sub(newStart) < minDuration {
		return dto.AuctionSessionDTO{}, ErrInvalidDate
	}

	// save to DB
	err = s.repo.Update(session)
	if err != nil {
		s.logger.Error("Failed to update auction session", "error", err)
		return dto.AuctionSessionDTO{}, ErrInvalidAuction
	}

	s.logger.Info("Auction session updated", "sessionID", id)
	return dto.AuctionSessionResponse(*session), nil
}

func (s *sessionService) Delete(id int64) error {
	session, err := s.repo.GetByID(id)
	if err != nil {
		return ErrSessionNotFoundID
	}

	// Convert DB times (stored as UTC) to local for comparison
	now := time.Now().In(wibLocation)
	sessionStart := session.StartTime.In(wibLocation)
	sessionEnd := session.EndTime.In(wibLocation)

	// cannot delete active (ongoing) session
	if sessionStart.Before(now) && sessionEnd.After(now) {
		return ErrActiveSession
	}

	// cannot delete expired session (keep historical data)
	if sessionEnd.Before(now) {
		return ErrExpiredSession
	}

	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete auction session", "error", err)
		return ErrInvalidAuction
	}

	s.logger.Info("Auction session deleted", "sessionID", id)
	return nil
}
