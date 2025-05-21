package models

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"size:191;uniqueIndex;not null"`
	Password   string `gorm:"size:255;not null"`
	Role       string `gorm:"size:50;default:user"`
	EmployeeID *uint
	Employee   *Employee `gorm:"foreignKey:EmployeeID"`
}
