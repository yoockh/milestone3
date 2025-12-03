package dto

import (
	"time"

	"milestone3/be/internal/entity"
)

type AuctionItemDTO struct {
	Title         string  `json:"title,omitempty" validate:"required"`
	Description   string  `json:"description,omitempty" validate:"required"`
	Category      string  `json:"category,omitempty" validate:"required"`
	Status        string  `json:"status,omitempty"`
	ID            int64   `json:"id,omitempty" validate:"omitempty"`
	UserID        int64   `json:"user_id,omitempty"`
	DonationID    int64   `json:"donation_id,omitempty" validate:"required"`
	SessionID     *int64  `json:"session_id,omitempty"`
	StartingPrice float64 `json:"starting_price,omitempty"`
	// Photos        []string  `json:"photos,omitempty" validate:"dive,url"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func AuctionItemRequest(d AuctionItemDTO) (entity.AuctionItem, error) {
	status := d.Status
	if status == "" {
		status = "scheduled"
	}

	return entity.AuctionItem{
		ID:            d.ID,
		DonationID:    d.DonationID,
		SessionID:     d.SessionID,
		Title:         d.Title,
		Description:   d.Description,
		Category:      d.Category,
		Status:        status,
		StartingPrice: d.StartingPrice,
		CreatedAt:     d.CreatedAt,
	}, nil
}

func AuctionItemResponse(m entity.AuctionItem) AuctionItemDTO {
	return AuctionItemDTO{
		ID:            m.ID,
		DonationID:    m.DonationID,
		SessionID:     m.SessionID,
		Title:         m.Title,
		Description:   m.Description,
		Category:      m.Category,
		Status:        m.Status,
		StartingPrice: m.StartingPrice,
		CreatedAt: m.CreatedAt,
	}
}

func AuctionItemResponses(ms []entity.AuctionItem) []AuctionItemDTO {
	res := make([]AuctionItemDTO, 0, len(ms))
	for _, m := range ms {
		res = append(res, AuctionItemResponse(m))
	}
	return res
}

type AuctionSessionDTO struct {
	Name      string    `json:"name,omitempty" validate:"required"`
	ID        int64     `json:"id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty" validate:"required"`
	EndTime   time.Time `json:"end_time,omitempty" validate:"required"`
}

func AuctionSessionResponse(m entity.AuctionSession) AuctionSessionDTO {
	return AuctionSessionDTO{
		Name:      m.Name,
		ID:        m.ID,
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
	}
}

func AuctionSessionResponses(ms []entity.AuctionSession) []AuctionSessionDTO {
	res := make([]AuctionSessionDTO, 0, len(ms))
	for _, m := range ms {
		res = append(res, AuctionSessionResponse(m))
	}
	return res
}

func AuctionSessionRequest(d AuctionSessionDTO) entity.AuctionSession {
	return entity.AuctionSession{
		Name:      d.Name,
		StartTime: d.StartTime,
		EndTime:   d.EndTime,
	}
}
