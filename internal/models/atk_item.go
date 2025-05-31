package models

import "time"

type ATKItem struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:191;not null" json:"name"`
	Category      string     `gorm:"size:100" json:"category"`
	Condition     string     `gorm:"size:100" json:"condition"`
	Quantity      int        `json:"quantity"`
	WarehouseID   uint       `json:"warehouse_id"`
	Warehouse     Warehouse  `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	ResponsibleID *uint      `json:"responsible_id,omitempty"`
	Responsible   *Employee  `gorm:"foreignKey:ResponsibleID;references:ID" json:"responsible,omitempty"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	Price         float64    `json:"price"`
	Value         float64    `json:"value"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
