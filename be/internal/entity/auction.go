package entity

import "time"

type AuctionSession struct {
	Name      string    `gorm:"size:255;not null" json:"name"`
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Items []AuctionItem `gorm:"foreignKey:SessionID;constraint:OnDelete:SET NULL" json:"items,omitempty"`
}

type AuctionItem struct {
	Title         string    `gorm:"size:255;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Category      string    `gorm:"size:255" json:"category"`
	Status        string    `gorm:"size:50;default:'pending';not null" json:"status"`
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DonationID    int64     `gorm:"not null" json:"donation_id"`
	StartingPrice float64   `gorm:"not null" json:"starting_price"`
	SessionID     *int64    `gorm:"null" json:"session_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	Session *AuctionSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Photos  []DonationPhoto `gorm:"foreignKey:DonationID;constraint:OnDelete:CASCADE" json:"photos,omitempty"`
}
