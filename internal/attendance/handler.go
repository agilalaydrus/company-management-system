package attendance

import (
	"fmt"
	"io"
	"metro-backend/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAttendanceHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeIDStr := c.PostForm("employee_id")
		latStr := c.PostForm("lat")
		longStr := c.PostForm("long")

		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
			return
		}
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
			return
		}
		long, err := strconv.ParseFloat(longStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
			return
		}

		// Cek sudah absen hari ini
		var existing models.Attendance
		today := time.Now().Format("2006-01-02")
		err = db.Where("employee_id = ? AND DATE(clock_in) = ?", employeeID, today).First(&existing).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "already clocked in today"})
			return
		}

		// Upload foto ClockIn
		header, err := c.FormFile("photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo is required"})
			return
		}
		file, err := header.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open photo"})
			return
		}
		defer file.Close()

		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
			return
		}

		// ✅ DI SINI: deklarasi filePath sebelum digunakan
		filename := fmt.Sprintf("%d_clockin_%s", time.Now().Unix(), filepath.Base(header.Filename))
		filePath := filepath.Join(uploadDir, filename)

		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo"})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write photo"})
			return
		}

		attendance := models.Attendance{
			EmployeeID:      uint(employeeID),
			Latitude:        lat,
			Longitude:       long,
			ClockIn:         time.Now(),
			ClockInPhotoURL: filePath, // pakai filePath yang sudah benar
		}

		if err := db.Create(&attendance).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save attendance"})
			return
		}

		var result models.Attendance
		if err := db.Preload("Employee").Preload("Employee.Company").First(&result, attendance.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendance with relations"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "clock in successful",
			"attendance": result,
		})
	}
}

func ClockOutHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeIDStr := c.PostForm("employee_id")
		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil || employeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
			return
		}

		// Ambil record attendance hari ini yang belum ClockOut
		var attendance models.Attendance
		today := time.Now().Format("2006-01-02")
		err = db.Where("employee_id = ? AND DATE(clock_in) = ? AND clock_out IS NULL", employeeID, today).
			First(&attendance).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "attendance not found or already clocked out"})
			return
		}

		// Upload foto ClockOut
		header, err := c.FormFile("photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo is required for clock out"})
			return
		}
		file, err := header.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open photo"})
			return
		}
		defer file.Close()

		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
			return
		}

		filename := fmt.Sprintf("%d_clockout_%s", time.Now().Unix(), filepath.Base(header.Filename))
		filePath := filepath.Join(uploadDir, filename)

		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo"})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write photo"})
			return
		}

		// Update ClockOut dan foto baru
		now := time.Now()
		attendance.ClockOut = &now
		attendance.ClockOutPhotoURL = filePath

		if err := db.Save(&attendance).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update clock out"})
			return
		}

		var result models.Attendance
		if err := db.Preload("Employee").Preload("Employee.Company").First(&result, attendance.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated attendance"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "clock out successful with photo",
			"attendance": result,
		})
	}
}

func GetAttendanceHistoryByEmployee(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeIDStr := c.Param("employee_id")
		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid employee_id"})
			return
		}

		var attendances []models.Attendance
		if err := db.Where("employee_id = ?", employeeID).Order("clock_in desc").Find(&attendances).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch attendance history"})
			return
		}

		c.JSON(200, attendances)
	}
}
