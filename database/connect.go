package database

import (
    "fmt"
    "os"
    // "payment/dbutils"
    "github.com/joho/godotenv"

    "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
    err := godotenv.Load()
    if err != nil {
        fmt.Println(err.Error())
    }

    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    database_name := os.Getenv("DB_DATABASE_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}