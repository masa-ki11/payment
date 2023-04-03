package models

import (
	"fmt"
	"payment/database"
	"time"
)

func SaveUser(user *User) error {
	db, err := database.Connect()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return err
	}

	result := db.Create(&user)

	if result.Error != nil {
		fmt.Println("SaveUser error:", result.Error)
		return result.Error
	}

	fmt.Println("Inserted user with ID:", user.ID)
	return nil
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"-"`
	Admin    bool
	DeleteFlag    bool `gorm:"column:delete_flag"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
