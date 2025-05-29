package models

import "time"

type Payslip struct {
	ID          uint      `json:"ID"`
	EmployeeID  uint      `json:"EmployeeID"`
	Period      string    `json:"Period"`
	BaseSalary  int       `json:"BaseSalary"`
	Allowance   int       `json:"Allowance"`
	Deduction   int       `json:"Deduction"`
	TotalSalary int       `json:"TotalSalary"`
	GeneratedAt time.Time `json:"GeneratedAt"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Employee Employee `gorm:"foreignKey:EmployeeID" json:"Employee"`
}
