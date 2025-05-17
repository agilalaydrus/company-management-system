package models

import "time"

type Attendance struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"not null"`
	PhotoURL   string    `gorm:"not null"`
	Lat        float64   `gorm:"not null"`
	Long       float64   `gorm:"not null"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}
