package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/models"
)

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.Company
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create company"})
			return
		}
		c.JSON(http.StatusCreated, input)
	}
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var companies []models.Company
		if err := db.Find(&companies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch companies"})
			return
		}
		c.JSON(http.StatusOK, companies)
	}
}

func Detail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var company models.Company
		if err := db.First(&company, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusOK, company)
	}
}

func Update(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var company models.Company
		if err := db.First(&company, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		if err := c.ShouldBindJSON(&company); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		db.Save(&company)
		c.JSON(http.StatusOK, company)
	}
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		db.Delete(&models.Company{}, id)
		c.Status(http.StatusNoContent)
	}
}
