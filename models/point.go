package models

import (
	// "fmt"
	// "payment/database"
	"time"
)

type Point struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	Point     uint `json:"point"`
	CreatedAt time.Time
	UpdatedAt time.Time
}