package database

import (
	"log"
	"os"

	"gostockly/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.Store{}, &models.Inventory{}, &models.StockGroup{}, &models.Company{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}
