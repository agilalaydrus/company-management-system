package db

import (
	"gorm.io/gorm"
	"log"
	"metro-backend/internal/models"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Company{},
		&models.Employee{},
		&models.Attendance{},
		&models.Leave{},
		&models.Letter{},
		&models.Payslip{},
		&models.User{},
		&models.LetterTemplate{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate tables: %v", err)
	} else {
		log.Println("✅ Successfully migrated tables")
	}
}
