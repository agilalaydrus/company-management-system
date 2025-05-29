package models

import "time"

type Attendance struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	EmployeeID       uint       `json:"employee_id"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	ClockIn          time.Time  `json:"clock_in"`
	ClockOut         *time.Time `json:"clock_out,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ClockInPhotoURL  string     `json:"clock_in_photo_url"`
	ClockOutPhotoURL string     `json:"clock_out_photo_url,omitempty"`
	Employee         Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
}
