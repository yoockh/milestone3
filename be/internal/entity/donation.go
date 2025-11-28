package entity

import (
	"time"
)

type Donation struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	// User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"` // Assuming User entity exists elsewhere
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:255" json:"category"`
	Condition   string    `gorm:"size:255" json:"condition"`
	Status      string    `gorm:"type:donation_status;default:'pending';not null" json:"status"` // enum: pending, verified_for_auction, verified_for_donation
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Photos []DonationPhoto `gorm:"foreignKey:DonationID;constraint:OnDelete:CASCADE" json:"photos,omitempty"`
}

type DonationPhoto struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DonationID uint   `gorm:"not null" json:"donation_id"`
	URL        string `gorm:"size:255" json:"url"`
}
