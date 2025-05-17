package models

import "time"

type Letter struct {
	ID           uint      `gorm:"primaryKey"`
	EmployeeID   uint      `gorm:"not null"`
	CompanyName  string    `gorm:"size:191;not null"`
	LetterNumber string    `gorm:"size:191;uniqueIndex;not null"`
	LetterType   string    `gorm:"size:100"`
	Content      string    `gorm:"type:text"`
	IssuedDate   time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
