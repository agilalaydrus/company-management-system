package letter

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/models"
)

// CreateLetter handles the creation of letters with predefined templates
func CreateLetter(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			EmployeeID uint   `json:"employee_id"`
			CompanyID  uint   `json:"company_id"`
			TemplateID uint   `json:"template_id"`
			Content    string `json:"content"`     // user bisa kosongkan jika mau pakai default
			IssuedDate string `json:"issued_date"` // format: YYYY-MM-DD
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		date, err := time.Parse("2006-01-02", req.IssuedDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}

		// Ambil template
		var tmpl models.LetterTemplate
		if err := db.First(&tmpl, req.TemplateID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "template not found"})
			return
		}

		// Baca isi file template
		defaultContent := ""
		tmplPath := filepath.Join("templates", tmpl.FilePath)
		raw, err := os.ReadFile(tmplPath)
		if err == nil {
			defaultContent = string(raw)
		}

		// Gunakan default content jika user tidak isi
		content := req.Content
		if content == "" {
			content = defaultContent
		}

		// Generate nomor surat
		var count int64
		db.Model(&models.Letter{}).Where("template_id = ?", req.TemplateID).Count(&count)
		number := fmt.Sprintf("%d/TEMP/%03d", time.Now().Year(), count+1)

		// Simpan surat
		letter := models.Letter{
			EmployeeID:   req.EmployeeID,
			CompanyID:    req.CompanyID,
			TemplateID:   req.TemplateID,
			LetterNumber: number,
			Content:      content,
			IssuedDate:   date,
		}

		if err := db.Create(&letter).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to create letter",
				"details": err.Error(),
			})
			return
		}

		if err := db.Preload("Employee").
			Preload("Company").
			Preload("Template").
			First(&letter, letter.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload letter with relations"})
			return
		}
		c.JSON(http.StatusCreated, letter)
	}
}

// GetLettersFiltered retrieves all letters with optional filtering
func GetLettersFiltered(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		template := c.Query("template")
		company := c.Query("company")

		var letters []models.Letter
		query := db.Preload("Employee").Preload("Company").Preload("Template").Order("issued_date desc")

		if template != "" {
			query = query.Joins("JOIN letter_templates ON letter_templates.id = letters.template_id").
				Where("letter_templates.name = ?", template)
		}

		if company != "" {
			query = query.Joins("JOIN companies ON companies.id = letters.company_id").
				Where("companies.name = ?", company)
		}

		if err := query.Find(&letters).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch letters"})
			return
		}

		c.JSON(http.StatusOK, letters)
	}
}

// GenerateLetterHTML renders a letter as HTML using the associated template
func GenerateLetterHTML(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var letter models.Letter
		if err := db.Preload("Employee").Preload("Company").Preload("Template").First(&letter, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "letter not found"})
			return
		}

		// Gunakan template umum yang akan menampilkan Content sebagai HTML
		tmplPath := filepath.Join("templates", "base_content_layout.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load template layout"})
			return
		}

		// Render isi surat (Content HTML)
		c.Writer.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, letter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render content"})
			return
		}
	}
}

func ExportLetterPDF(db *gorm.DB, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid letter ID"})
			return
		}

		var letter models.Letter
		if err := db.Preload("Employee").Preload("Template").First(&letter, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "letter not found"})
			return
		}

		// Buat nama file: Surat_Template_Employee.pdf
		sanitize := func(s string) string {
			// Remove non-alphanumerics and replace space with _
			s = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(s, "")
			s = strings.ReplaceAll(s, " ", "_")
			return s
		}
		templateName := sanitize(letter.Template.Name)
		employeeName := sanitize(letter.Employee.Name)
		filename := fmt.Sprintf("Surat_%s_%s.pdf", templateName, employeeName)

		// Temp path & HTML URL
		outputPath := filepath.Join("/tmp", filename)
		htmlURL := fmt.Sprintf("%s/letters/%d/html", baseURL, id)

		cmd := exec.Command("wkhtmltopdf", htmlURL, outputPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to generate PDF",
				"details": string(output),
			})
			return
		}

		// Kirim file PDF
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		file, err := os.ReadFile(outputPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read PDF"})
			return
		}
		c.Writer.Write(file)
	}
}
