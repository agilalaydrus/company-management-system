package models

import "time"

type InventoryItem struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:191;not null" json:"name"`
	Category      string     `gorm:"size:100" json:"category"`
	Condition     string     `gorm:"size:100" json:"condition"`
	Quantity      int        `json:"quantity"`
	WarehouseID   uint       `json:"warehouse_id"`
	Warehouse     Warehouse  `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	ResponsibleID *uint      `json:"responsible_id,omitempty"`
	Responsible   *Employee  `gorm:"foreignKey:ResponsibleID" json:"responsible,omitempty"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	Value         float64    `json:"value,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
