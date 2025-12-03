package entity

import "time"

type Bid struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	SessionID int64     `gorm:"not null;index"`
	ItemID    int64     `gorm:"not null;index"`
	Amount    float64   `gorm:"not null"`
	UserID    int64     `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
