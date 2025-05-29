package models

import "time"

type Company struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Name       string      `gorm:"size:191;not null" json:"name"`
	Address    string      `gorm:"size:255" json:"address"`
	Employees  []Employee  `json:"employees,omitempty"`
	Letters    []Letter    `json:"letters,omitempty"`
	Warehouses []Warehouse `json:"warehouses,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
