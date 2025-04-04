package db

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error loading env ", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var er error
	db, er := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if er != nil {
		log.Fatal("Failed to connect to database:", er)
	}

	DB = db
	fmt.Println("connected to PostgreSQL ")

	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Message{}, &models.Like{}, &models.Comments{})

}

func CreateUsers(username, email, password string) error {
	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	result := DB.Create(&user)

	return result.Error
}
