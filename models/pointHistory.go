package models

import (
	// "fmt"
	// "payment/database"
	"time"
)

type PointHistory struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	Point     uint `json:"point"`
	Action      string `json:"action"`
	Details      string `json:"details"`
	CreatedAt time.Time
	CreatedAtFormatted string `json:"created_at_formatted" gorm:"-"`
}