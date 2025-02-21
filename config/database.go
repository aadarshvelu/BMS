package config

import (
	"fmt"
	"log"

	"github.com/aadarshvelu/bms/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "Test@123"),
		GetEnv("DB_NAME", "bms"),
		GetEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the Book model
	err = db.AutoMigrate(&models.Book{})

	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	// assign to global pointer
	DB = db
}
