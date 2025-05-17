package models

import "time"

type Payslip struct {
	ID          uint   `gorm:"primaryKey"`
	EmployeeID  uint   `gorm:"not null"`
	Period      string `gorm:"size:20"` // format: 2024-05
	BaseSalary  int    `gorm:"not null"`
	Allowance   int
	Deduction   int
	TotalSalary int
	GeneratedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
