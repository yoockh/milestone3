package scheduler

import (
	"log/slog"
	"milestone3/be/internal/service"
	"time"

	"github.com/go-co-op/gocron"
)

type BidScheduler struct {
	svc    service.BidService
	logger *slog.Logger
}

func NewBidScheduler(bidService service.BidService, logger *slog.Logger) *BidScheduler {
	return &BidScheduler{
		svc:    bidService,
		logger: logger,
	}
}

func (s *BidScheduler) Start() {
	// set as local time
	scheduler := gocron.NewScheduler(time.Local)

	// save expired sessions to DB every 1 minute
	_, err := scheduler.Every(1).Minute().Do(func() {
		s.logger.Info("Running expired sessions cleanup...")
		if err := s.svc.SaveKeyToDB(); err != nil {
			s.logger.Error("Failed to save expired sessions", "error", err)
		}
	})

	if err != nil {
		s.logger.Error("Failed to schedule bid sync", "error", err)
		return
	}

	// delete key value at 12 AM daily
	_, err = scheduler.Every(1).Day().At("00:00").Do(func() {
		s.logger.Info("Running midnight Redis cleanup...")
		if err = s.svc.DeleteKeyValue(); err != nil {
			s.logger.Error("Failed to cleanup Redis at midnight", "error", err)
		}
	})

	if err != nil {
		s.logger.Error("Failed to schedule midnight cleanup", "error", err)
		return
	}

	scheduler.StartAsync()
	s.logger.Info("Bid scheduler started")
	s.logger.Info("- Sync to DB: every 1 minute")
	s.logger.Info("- Redis cleanup: daily at 00:00")
}
