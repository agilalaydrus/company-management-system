package attendance

import (
	"fmt"
	"io"
	"metro-backend/internal/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAttendanceHandlers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeIDStr := c.PostForm("employee_id")
		latStr := c.PostForm("lat")
		longStr := c.PostForm("long")

		employeeID, _ := strconv.Atoi(employeeIDStr)
		lat, _ := strconv.ParseFloat(latStr, 64)
		long, _ := strconv.ParseFloat(longStr, 64)

		file, header, err := c.Request.FormFile("photo")
		if err != nil {
			c.JSON(400, gin.H{"error": "photo is required"})
			return
		}
		defer file.Close()

		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filepath.Base(header.Filename))
		os.MkdirAll("uploads", os.ModePerm)
		out, _ := os.Create(filename)
		io.Copy(out, file)

		attendance := models.Attendance{
			EmployeeID: uint(employeeID),
			Lat:        lat,
			Long:       long,
			PhotoURL:   filename,
		}

		if err := db.Create(&attendance).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to save attendance"})
			return
		}

		c.JSON(201, gin.H{"message": "attendance recorded", "id": attendance.ID})
	}
}
