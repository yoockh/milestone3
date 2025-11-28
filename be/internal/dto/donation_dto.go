package dto

import (
	"time"

	"milestone3/be/internal/entity"
)

type DonationDTO struct {
	ID          uint      `json:"id,omitempty" validate:"omitempty"`
	UserID      uint      `json:"user_id,omitempty" validate:"required"`
	Title       string    `json:"title,omitempty" validate:"required"`
	Description string    `json:"description,omitempty" validate:"required"`
	Category    string    `json:"category,omitempty" validate:"required"`
	Condition   string    `json:"condition,omitempty" validate:"required"`
	Status      string    `json:"status,omitempty" validate:"required"`
	Photos      []string  `json:"photos,omitempty" validate:"dive,url"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// DonationRequest converts DTO to entity.Donation
func DonationRequest(d DonationDTO) (entity.Donation, error) {
	photos := make([]entity.DonationPhoto, 0, len(d.Photos))
	for _, u := range d.Photos {
		photos = append(photos, entity.DonationPhoto{URL: u})
	}
	return entity.Donation{
		ID:          d.ID,
		UserID:      d.UserID,
		Title:       d.Title,
		Description: d.Description,
		Category:    d.Category,
		Condition:   d.Condition,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		Photos:      photos,
	}, nil
}

// DonationResponse converts entity.Donation to DTO
func DonationResponse(m entity.Donation) DonationDTO {
	photos := make([]string, 0, len(m.Photos))
	for _, p := range m.Photos {
		photos = append(photos, p.URL)
	}
	return DonationDTO{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		Category:    m.Category,
		Condition:   m.Condition,
		Status:      m.Status,
		Photos:      photos,
		CreatedAt:   m.CreatedAt,
	}
}

// DonationResponses converts slice of entity.Donation to slice of DTOs
func DonationResponses(ms []entity.Donation) []DonationDTO {
	res := make([]DonationDTO, 0, len(ms))
	for _, m := range ms {
		res = append(res, DonationResponse(m))
	}
	return res
}

type DonationPhoto struct {
	URL string
}

type Donation struct {
	ID          uint
	UserID      uint
	Title       string
	Description string
	Category    string
	Condition   string
	Status      string
	Photos      []DonationPhoto
	CreatedAt   time.Time
}
