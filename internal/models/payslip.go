package models

import "time"

type Payslip struct {
	ID          uint   `gorm:"primaryKey"`
	EmployeeID  uint   `gorm:"not null"`
	Period      string `gorm:"size:20"` // format: YYYY-MM
	BaseSalary  int    `gorm:"not null"`
	Allowance   int
	Deduction   int
	TotalSalary int
	GeneratedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Employee    Employee `gorm:"foreignKey:EmployeeID"`
}
