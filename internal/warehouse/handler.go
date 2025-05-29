package warehouse

import (
	"net/http"
	"strconv"

	"metro-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create Warehouse
func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name      string `json:"name" binding:"required"`
			Location  string `json:"location"`
			CompanyID uint   `json:"company_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		warehouse := models.Warehouse{
			Name:      input.Name,
			Location:  input.Location,
			CompanyID: input.CompanyID,
		}

		if err := db.Create(&warehouse).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create warehouse"})
			return
		}

		// Reload warehouse dengan preload company
		if err := db.Preload("Company").First(&warehouse, warehouse.ID).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to load warehouse with company"})
			return
		}

		c.JSON(201, warehouse)
	}
}

// List Warehouses (with preload Company)
func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var warehouses []models.Warehouse
		if err := db.Preload("Company").Find(&warehouses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list warehouses"})
			return
		}
		c.JSON(http.StatusOK, warehouses)
	}
}

// Detail Warehouse (with preload Company, InventoryItems, ATKItems)
func Detail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
			return
		}

		var warehouse models.Warehouse
		if err := db.Preload("Company").
			Preload("InventoryItems").
			Preload("ATKItems").
			First(&warehouse, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
			return
		}

		c.JSON(http.StatusOK, warehouse)
	}
}

// Update Warehouse
func Update(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
			return
		}

		var warehouse models.Warehouse
		if err := db.First(&warehouse, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
			return
		}

		var input struct {
			Name     string `json:"name"`
			Location string `json:"location"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name != "" {
			warehouse.Name = input.Name
		}
		if input.Location != "" {
			warehouse.Location = input.Location
		}

		if err := db.Save(&warehouse).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update warehouse"})
			return
		}

		c.JSON(http.StatusOK, warehouse)
	}
}

// Delete Warehouse
func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
			return
		}

		if err := db.Delete(&models.Warehouse{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete warehouse"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
