package letter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"html/template"
	"metro-backend/internal/models"
	"net/http"
	"path/filepath"
	"time"
)

func CreateLetter(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			EmployeeID  uint   `json:"employee_id"`
			CompanyName string `json:"company_name"`
			LetterType  string `json:"letter_type"`
			Content     string `json:"content"`
			IssuedDate  string `json:"issued_date"` // format: YYYY-MM-DD
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		date, err := time.Parse("2006-01-02", req.IssuedDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		// generate nomor surat otomatis
		var count int64
		db.Model(&models.Letter{}).Where("letter_type = ?", req.LetterType).Count(&count)
		number := fmt.Sprintf("%d/%s/%03d", time.Now().Year(), req.LetterType, count+1)

		letter := models.Letter{
			EmployeeID:   req.EmployeeID,
			CompanyName:  req.CompanyName,
			LetterNumber: number,
			LetterType:   req.LetterType,
			Content:      req.Content,
			IssuedDate:   date,
		}

		if err := db.Create(&letter).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create letter"})
			return
		}

		c.JSON(http.StatusCreated, letter)
	}
}

func GetLettersFiltered(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		letterType := c.Query("type")
		company := c.Query("company")

		var letters []models.Letter
		query := db.Order("issued_date desc")

		if letterType != "" {
			query = query.Where("letter_type = ?", letterType)
		}
		if company != "" {
			query = query.Where("company_name = ?", company)
		}

		if err := query.Find(&letters).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch letters"})
			return
		}

		c.JSON(200, letters)
	}
}

func GenerateLetterHTML(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var letter models.Letter
		if err := db.First(&letter, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "letter not found"})
			return
		}

		tmplPath := filepath.Join("templates", "letter_template.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			c.JSON(500, gin.H{"error": "template error"})
			return
		}

		c.Writer.Header().Set("Content-Type", "text/html")
		tmpl.Execute(c.Writer, letter)
	}
}
