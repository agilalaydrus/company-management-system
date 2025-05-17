package models

import "time"

type Employee struct {
	ID             uint       `gorm:"primaryKey"`
	Name           string     `gorm:"size:191;not null"`
	NIK            string     `gorm:"size:100;uniqueIndex;not null"`
	Position       string     `gorm:"size:100"`
	JoinDate       *time.Time `json:"join_date"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RemainingLeave int `gorm:"default:12"`
}
