package inventory

import (
	"net/http"
	"strconv"
	"time"

	"metro-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name          string  `json:"name" binding:"required"`
			Category      string  `json:"category"`
			Condition     string  `json:"condition"`
			Quantity      int     `json:"quantity" binding:"required"`
			WarehouseID   uint    `json:"warehouse_id" binding:"required"`
			ResponsibleID *uint   `json:"responsible_id"`
			PurchaseDate  *string `json:"purchase_date"` // format "2006-01-02"
			Value         float64 `json:"value"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var purchaseDate *time.Time
		if input.PurchaseDate != nil {
			t, err := time.Parse("2006-01-02", *input.PurchaseDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase_date format"})
				return
			}
			purchaseDate = &t
		}

		item := models.InventoryItem{
			Name:          input.Name,
			Category:      input.Category,
			Condition:     input.Condition,
			Quantity:      input.Quantity,
			WarehouseID:   input.WarehouseID,
			ResponsibleID: input.ResponsibleID,
			PurchaseDate:  purchaseDate,
			Value:         input.Value,
		}

		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create inventory item"})
			return
		}

		// Fetch ulang item dengan preload relasi
		if err := db.Preload("Warehouse").Preload("Responsible").First(&item, item.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load inventory item with relations"})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.InventoryItem
		if err := db.Preload("Warehouse").Preload("Responsible").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list inventory items"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func Detail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory item id"})
			return
		}

		var item models.InventoryItem
		if err := db.Preload("Warehouse").Preload("Responsible").First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "inventory item not found"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func Update(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory item id"})
			return
		}

		var item models.InventoryItem
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "inventory item not found"})
			return
		}

		var input struct {
			Name          string  `json:"name"`
			Category      string  `json:"category"`
			Condition     string  `json:"condition"`
			Quantity      int     `json:"quantity"`
			WarehouseID   uint    `json:"warehouse_id"`
			ResponsibleID *uint   `json:"responsible_id"`
			PurchaseDate  *string `json:"purchase_date"`
			Value         float64 `json:"value"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name != "" {
			item.Name = input.Name
		}
		if input.Category != "" {
			item.Category = input.Category
		}
		if input.Condition != "" {
			item.Condition = input.Condition
		}
		if input.Quantity != 0 {
			item.Quantity = input.Quantity
		}
		if input.WarehouseID != 0 {
			item.WarehouseID = input.WarehouseID
		}
		item.ResponsibleID = input.ResponsibleID

		if input.PurchaseDate != nil {
			t, err := time.Parse("2006-01-02", *input.PurchaseDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase_date format"})
				return
			}
			item.PurchaseDate = &t
		}
		item.Value = input.Value

		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory item"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory item id"})
			return
		}

		if err := db.Delete(&models.InventoryItem{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete inventory item"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
