package models

import "time"

type Employee struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Name           string     `gorm:"size:191;not null" json:"name"`
	NIK            string     `gorm:"size:100;uniqueIndex;not null" json:"nik"`
	Position       string     `gorm:"size:100" json:"position"`
	JoinDate       *time.Time `json:"join_date"`
	RemainingLeave int        `gorm:"default:12" json:"remaining_leave"`
	CompanyID      uint       `json:"company_id"`
	Company        Company    `gorm:"foreignKey:CompanyID" json:"company"`

	Attendances []Attendance `gorm:"foreignKey:EmployeeID" json:"attendances,omitempty"`
	Leaves      []Leave      `gorm:"foreignKey:EmployeeID" json:"leaves,omitempty"`
	Letters     []Letter     `gorm:"foreignKey:EmployeeID" json:"letters,omitempty"`
	Payslips    []Payslip    `gorm:"foreignKey:EmployeeID" json:"payslips,omitempty"`

	InventoryItems []InventoryItem `gorm:"foreignKey:ResponsibleID" json:"inventory_items,omitempty"`
	ATKItems       []ATKItem       `gorm:"foreignKey:ResponsibleID" json:"atk_items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
