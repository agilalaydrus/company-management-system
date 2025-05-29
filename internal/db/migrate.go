package db

import (
	"log"
	"metro-backend/internal/models"

	"gorm.io/gorm"
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
		&models.Warehouse{},
		&models.InventoryItem{},
		&models.ATKItem{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate tables: %v", err)
	} else {
		log.Println("✅ Successfully migrated tables")
	}
}
