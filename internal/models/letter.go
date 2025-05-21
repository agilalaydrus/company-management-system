package models

import "time"

type Letter struct {
	ID           uint      `gorm:"primaryKey"`
	EmployeeID   uint      `gorm:"not null"`
	CompanyID    uint      `gorm:"not null"`
	TemplateID   uint      `gorm:"not null"`
	LetterNumber string    `gorm:"size:191;uniqueIndex;not null"`
	Content      string    `gorm:"type:text"` // berisi isi surat hasil edit
	IssuedDate   time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Employee Employee       `gorm:"foreignKey:EmployeeID"`
	Company  Company        `gorm:"foreignKey:CompanyID"`
	Template LetterTemplate `gorm:"foreignKey:TemplateID"`
}
