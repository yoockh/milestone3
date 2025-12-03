package entity

import "time"

type Bid struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID    int64     `gorm:"column:auction_item_id;not null;index" json:"auction_item_id"`
	Amount    float64   `gorm:"not null" json:"amount"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Bid) TableName() string {
	return "bids"
}
