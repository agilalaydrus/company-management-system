package models

import "time"

type Company struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:191;not null"`
	Address   string `gorm:"size:255"`
	Employees []Employee
	Letters   []Letter
	CreatedAt time.Time
	UpdatedAt time.Time
}
