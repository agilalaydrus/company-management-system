package leave

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/models"
	"time"
)

func ApplyLeave(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			EmployeeID uint   `json:"employee_id"`
			StartDate  string `json:"start_date"` // format: YYYY-MM-DD
			EndDate    string `json:"end_date"`
			Reason     string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		start, err1 := time.Parse("2006-01-02", req.StartDate)
		end, err2 := time.Parse("2006-01-02", req.EndDate)
		if err1 != nil || err2 != nil || end.Before(start) {
			c.JSON(400, gin.H{"error": "invalid date range"})
			return
		}

		leave := models.Leave{
			EmployeeID: req.EmployeeID,
			StartDate:  start,
			EndDate:    end,
			Reason:     req.Reason,
			Status:     "pending",
		}

		if err := db.Create(&leave).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to apply leave"})
			return
		}

		c.JSON(201, gin.H{
			"message": "leave requested successfully",
			"id":      leave.ID,
		})
	}
}

func ApproveLeave(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var leave models.Leave
		if err := db.First(&leave, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "leave not found"})
			return
		}

		var emp models.Employee
		if err := db.First(&emp, leave.EmployeeID).Error; err != nil {
			c.JSON(404, gin.H{"error": "employee not found"})
			return
		}

		// hitung hari
		duration := int(leave.EndDate.Sub(leave.StartDate).Hours()/24) + 1
		if emp.RemainingLeave < duration {
			c.JSON(400, gin.H{"error": "not enough leave balance"})
			return
		}

		leave.Status = "approved"
		emp.RemainingLeave -= duration

		db.Save(&leave)
		db.Save(&emp)

		c.JSON(200, gin.H{"message": "leave approved"})
	}
}

func RejectLeave(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var leave models.Leave
		if err := db.First(&leave, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "leave not found"})
			return
		}
		leave.Status = "rejected"
		db.Save(&leave)

		c.JSON(200, gin.H{"message": "leave rejected"})
	}
}

func ListLeave(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var leave []models.Leave
		db.Order("created_at desc").Find(&leave)
		c.JSON(200, leave)
	}
}
