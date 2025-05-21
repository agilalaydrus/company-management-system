package models

import "time"

type Leave struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	Reason     string    `gorm:"type:text"`
	Status     string    `gorm:"default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Employee   Employee `gorm:"foreignKey:EmployeeID"`
}
