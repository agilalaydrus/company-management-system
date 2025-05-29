package models

import "time"

type Warehouse struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"size:191;not null" json:"name"`
	Location  string  `gorm:"size:255" json:"location"`
	CompanyID uint    `json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`

	InventoryItems []InventoryItem `gorm:"foreignKey:WarehouseID" json:"inventory_items,omitempty"`
	ATKItems       []ATKItem       `gorm:"foreignKey:WarehouseID" json:"atk_items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
