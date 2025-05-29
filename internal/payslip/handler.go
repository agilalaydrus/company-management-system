package payslip

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/models"
	"time"
)

func CreatePayslip(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			EmployeeID uint   `json:"employee_id"`
			Period     string `json:"period"` // e.g. 2024-05
			BaseSalary int    `json:"base_salary"`
			Allowance  int    `json:"allowance"`
			Deduction  int    `json:"deduction"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		total := req.BaseSalary + req.Allowance - req.Deduction

		slip := models.Payslip{
			EmployeeID:  req.EmployeeID,
			Period:      req.Period,
			BaseSalary:  req.BaseSalary,
			Allowance:   req.Allowance,
			Deduction:   req.Deduction,
			TotalSalary: total,
			GeneratedAt: time.Now(),
		}

		// Setelah menyimpan slip
		if err := db.Create(&slip).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create payslip"})
			return
		}

		// Ambil ulang slip dengan relasi employee dan company
		var enriched models.Payslip
		if err := db.
			Preload("Employee").
			Preload("Employee.Company").
			First(&enriched, slip.ID).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to retrieve enriched payslip"})
			return
		}

		c.JSON(201, enriched)

	}
}

func GetPayslips(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payslips []models.Payslip
		db.
			Preload("Employee").
			Preload("Employee.Company").
			Order("generated_at desc").
			Find(&payslips)
		c.JSON(200, payslips)
	}
}
