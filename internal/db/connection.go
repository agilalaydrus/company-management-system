package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("✅ Database connected.")
			return db, nil
		}
		fmt.Printf("⏳ Attempt %d: waiting for DB... (%v)\n", i, err)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("❌ DB connection failed after retries: %v", err)
}
