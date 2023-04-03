package models

import (
	"time"
)

type PasswordReset struct {
	ID    uint
	Email string
	Token string
	CreatedAt time.Time
	UpdatedAt time.Time
}