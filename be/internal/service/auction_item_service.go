package service

import (
	"log/slog"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"
	"time"
)

type itemsService struct {
	repo   repository.AuctionItemRepository
	logger *slog.Logger
	ai     repository.AIRepository
}

type AuctionItemService interface {
	Create(item *dto.AuctionItemDTO) (dto.AuctionItemDTO, error)
	GetAll() ([]dto.AuctionItemDTO, error)
	GetByID(id int64) (dto.AuctionItemDTO, error)
	Update(id int64, item *dto.AuctionItemUpdateDTO) (dto.AuctionItemDTO, error)
	Delete(id int64) error
	CheckAndStartScheduledItems() error
}

func NewAuctionItemService(r repository.AuctionItemRepository, aiRepo repository.AIRepository, logger *slog.Logger) AuctionItemService {
	return &itemsService{repo: r, logger: logger, ai: aiRepo}
}

const DefaultStartingPrice = 10000

func (s *itemsService) Create(itemDTO *dto.AuctionItemDTO) (dto.AuctionItemDTO, error) {
	item, err := dto.AuctionItemRequest(*itemDTO)
	if err != nil {
		return dto.AuctionItemDTO{}, ErrInvalidAuction
	}

	estimationReq := repository.PriceEstimationRequest{
		Name:        itemDTO.Title,
		Category:    itemDTO.Category,
		Description: itemDTO.Description,
	}

	estimatedPrice, err := s.ai.EstimateStartingPrice(estimationReq)
	if err != nil {
		s.logger.Warn("EstimateStartingPrice failed, falling back to provided/default price", "error", err)
		if itemDTO.StartingPrice > 0 {
			estimatedPrice = itemDTO.StartingPrice
		} else {
			estimatedPrice = DefaultStartingPrice
		}
	}

	item.StartingPrice = estimatedPrice

	if item.Status == "" {
		item.Status = "scheduled"
	}

	err = s.repo.Create(&item)
	if err != nil {
		s.logger.Error("Failed to create auction item", "error", err)
		return dto.AuctionItemDTO{}, ErrInvalidAuction
	}

	return dto.AuctionItemResponse(item), nil
}

func (s *itemsService) GetAll() ([]dto.AuctionItemDTO, error) {
	items, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all auction items", "error", err)
		return nil, ErrAuctionNotFound
	}

	var itemDTOs []dto.AuctionItemDTO
	for _, item := range items {
		itemDTOs = append(itemDTOs, dto.AuctionItemResponse(item))
	}

	return itemDTOs, nil
}

func (s *itemsService) GetByID(id int64) (dto.AuctionItemDTO, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return dto.AuctionItemDTO{}, ErrAuctionNotFoundID
	}
	return dto.AuctionItemResponse(*item), nil
}

func (s *itemsService) Update(id int64, updateDTO *dto.AuctionItemUpdateDTO) (dto.AuctionItemDTO, error) {
	existingItem, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction item by ID for update", "error", err)
		return dto.AuctionItemDTO{}, ErrAuctionNotFoundID
	}

	if existingItem.Status == "finished" {
		s.logger.Warn("Cannot update finished auction item", "itemID", id, "status", existingItem.Status)
		return dto.AuctionItemDTO{}, ErrAuctionFinished
	}

	// validate and apply updates
	if updateDTO.Title != nil {
		existingItem.Title = *updateDTO.Title
	}
	if updateDTO.Description != nil {
		existingItem.Description = *updateDTO.Description
	}
	if updateDTO.Category != nil {
		existingItem.Category = *updateDTO.Category
	}
	if updateDTO.StartingPrice != nil {
		if *updateDTO.StartingPrice < 0 {
			s.logger.Warn("Invalid starting price", "price", *updateDTO.StartingPrice)
			return dto.AuctionItemDTO{}, ErrInvalidAuction
		}
		existingItem.StartingPrice = *updateDTO.StartingPrice
	}
	if updateDTO.Status != nil {
		newStatus := *updateDTO.Status
		// status transition rules
		switch existingItem.Status {
		case "scheduled":
			// to ongoing
			if newStatus != "ongoing" && newStatus != "scheduled" {
				s.logger.Warn("Invalid status transition", "from", existingItem.Status, "to", newStatus)
				return dto.AuctionItemDTO{}, ErrInvalidAuction
			}
		case "ongoing":
			// to finished
			if newStatus != "finished" && newStatus != "ongoing" {
				s.logger.Warn("Invalid status transition", "from", existingItem.Status, "to", newStatus)
				return dto.AuctionItemDTO{}, ErrInvalidAuction
			}
		case "finished":
			// cannot change from finished
			s.logger.Warn("Cannot change status from finished", "itemID", id)
			return dto.AuctionItemDTO{}, ErrAuctionFinished
		}
		existingItem.Status = newStatus
	}
	if updateDTO.SessionID != nil {
		// cannot change session if auction is ongoing
		if existingItem.Status == "ongoing" {
			s.logger.Warn("Cannot change session for ongoing auction", "itemID", id)
			return dto.AuctionItemDTO{}, ErrActiveSession
		}
		existingItem.SessionID = updateDTO.SessionID
	}
	if updateDTO.DonationID != nil {
		existingItem.DonationID = *updateDTO.DonationID
	}

	err = s.repo.Update(existingItem)
	if err != nil {
		s.logger.Error("Failed to update auction item", "error", err)
		return dto.AuctionItemDTO{}, ErrInvalidAuction
	}

	return dto.AuctionItemResponse(*existingItem), nil
}

func (s *itemsService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrAuctionNotFoundID
	}

	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete auction item", "error", err)
		return ErrInvalidAuction
	}
	return nil
}

// CheckAndStartScheduledItems (for cron job) to start if start_time has passed
func (s *itemsService) CheckAndStartScheduledItems() error {
	items, err := s.repo.GetScheduledItems()
	if err != nil {
		return err
	}

	// Convert both to same timezone for comparison
	now := time.Now().In(wibLocation)
	updatedCount := 0

	for _, item := range items {
		// check if item has a session
		if item.Session == nil {
			continue
		}

		// DB stores UTC, convert to WIB for comparison
		sessionStart := item.Session.StartTime.In(wibLocation)

		// check if current time >= session start_time
		if !now.Before(sessionStart) {
			// Time to start this auction
			item.Status = "ongoing"
			if err := s.repo.Update(&item); err != nil {
				s.logger.Error("Failed to start auction item", "itemID", item.ID, "error", err)
				continue
			}
			updatedCount++
		}
	}

	if updatedCount > 0 {
		s.logger.Info("Auto-started auction items", "count", updatedCount)
	}

	return nil
}
