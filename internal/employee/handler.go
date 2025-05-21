package employee

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/models"
	"net/http"
	"strconv"
)

func CreateEmployee(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var emp models.Employee
		if err := c.ShouldBindJSON(&emp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if emp.JoinDate == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "join_date is required"})
			return
		}

		// Set default leave balance
		emp.RemainingLeave = 12

		if err := db.Create(&emp).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee"})
			return
		}
		c.JSON(http.StatusCreated, emp)
	}
}

func GetAllEmployees(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var employees []models.Employee
		db.Find(&employees)
		c.JSON(http.StatusOK, employees)
	}
}

func GetEmployeeByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var emp models.Employee
		if err := db.First(&emp, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, emp)
	}
}

func UpdateEmployee(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var emp models.Employee
		if err := db.First(&emp, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if err := c.ShouldBindJSON(&emp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&emp)
		c.JSON(http.StatusOK, emp)
	}
}

func DeleteEmployee(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
			return
		}

		var employee models.Employee
		if err := db.First(&employee, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}

		if err := db.Delete(&employee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete employee"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "employee deleted successfully",
			"id":      id,
		})
	}
}
