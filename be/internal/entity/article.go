package entity

import (
	"time"
)

type Article struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Week      int       `json:"week"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
