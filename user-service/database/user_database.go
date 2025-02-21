package database

import (
	"fmt"
	"log"
	"os"

	"user-service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable that holds the database connection
// var DB *sql.DB   // for raw SQL queries

// using Go ORM
var DB *gorm.DB

func ConnectDB() {

	er := godotenv.Load()

	if er != nil {
		log.Fatal("error loading env ")
	}

	//url for DB from .env
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// To Check if the database is reachable
	// err = DB.Ping()
	// if err != nil {
	// 	log.Fatal("Database ping failed:", err)
	// }

	DB = db

	fmt.Println("connected to PostgreSQL ")

	//migration

	DB.AutoMigrate(&models.User{})
}
