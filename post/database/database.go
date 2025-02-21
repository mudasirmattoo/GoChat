package database

import (
	"fmt"
	"log"
	"os"

	"post-service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database credentials from environment variables
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Successfully connected to PostgreSQL!")

	// Verify connection using the underlying sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get SQL database instance:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	// Migrate the database schema
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	fmt.Println("Database migration successful!")
}
