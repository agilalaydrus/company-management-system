package db

import (
	"gorm.io/gorm"
	"log"
	"metro-backend/internal/models"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Attendance{},
		&models.Employee{},
		&models.Leave{},
		&models.Letter{},
		&models.Payslip{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate tables: %v", err)
	} else {
		log.Println("✅ Successfully migrated tables")
	}
}
