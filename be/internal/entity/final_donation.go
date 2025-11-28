package entity

import (
	"time"
)

type FinalDonation struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DonationID uint      `gorm:"not null" json:"donation_id"`
	Notes      string    `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
