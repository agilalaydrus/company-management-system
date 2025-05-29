package models

import "time"

type ATKItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:191;not null" json:"name"`
	Unit          string    `gorm:"size:50" json:"unit"`
	Quantity      int       `json:"quantity"`
	WarehouseID   uint      `json:"warehouse_id"`
	Warehouse     Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	ResponsibleID *uint     `json:"responsible_id,omitempty"`
	Responsible   *Employee `gorm:"foreignKey:ResponsibleID" json:"responsible,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
