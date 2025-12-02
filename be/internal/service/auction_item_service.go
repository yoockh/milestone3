package service

import (
	"log/slog"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/repository"
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
	Update(id int64, item *dto.AuctionItemDTO) (dto.AuctionItemDTO, error)
	Delete(id int64) error
}

func NewAuctionItemService(r repository.AuctionItemRepository, aiRepo repository.AIRepository, logger *slog.Logger) AuctionItemService {
	return &itemsService{repo: r, logger: logger, ai: aiRepo}
}

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
		// AI estimation failed — log and fallback to provided starting price or sensible default
		s.logger.Warn("EstimateStartingPrice failed, falling back to provided/default price", "error", err)
		// prefer client-provided StartingPrice if present
		if itemDTO.StartingPrice > 0 {
			estimatedPrice = itemDTO.StartingPrice
		} else {
			// sensible default to avoid blocking creation
			estimatedPrice = 100
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
	{
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
}

func (s *itemsService) GetByID(id int64) (dto.AuctionItemDTO, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return dto.AuctionItemDTO{}, ErrAuctionNotFoundID
	}
	return dto.AuctionItemResponse(*item), nil
}

func (s *itemsService) Update(id int64, itemDTO *dto.AuctionItemDTO) (dto.AuctionItemDTO, error) {
	existingItem, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get auction item by ID for update", "error", err)
		return dto.AuctionItemDTO{}, ErrAuctionNotFoundID
	}

	updatedItem, err := dto.AuctionItemRequest(*itemDTO)
	if err != nil {
		return dto.AuctionItemDTO{}, ErrInvalidAuction
	}

	updatedItem.ID = existingItem.ID
	updatedItem.CreatedAt = existingItem.CreatedAt

	err = s.repo.Update(&updatedItem)
	if err != nil {
		s.logger.Error("Failed to update auction item", "error", err)
		return dto.AuctionItemDTO{}, ErrInvalidAuction
	}

	return dto.AuctionItemResponse(updatedItem), nil
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
