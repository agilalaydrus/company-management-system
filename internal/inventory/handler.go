package inventory

import (
	"metro-backend/internal/models"
	"net/http"
	"strconv"
	"time"

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
			Price         float64 `json:"price" binding:"required"`
			WarehouseID   uint    `json:"warehouse_id" binding:"required"`
			ResponsibleID *uint   `json:"responsible_id"`
			PurchaseDate  *string `json:"purchase_date"`
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

		// Check existing inventory item
		var existing models.InventoryItem
		err := db.Where("name = ? AND category = ? AND warehouse_id = ? AND responsible_id = ?",
			input.Name, input.Category, input.WarehouseID, input.ResponsibleID).
			First(&existing).Error

		if err == nil {
			existing.Quantity += input.Quantity
			existing.Price = input.Price
			existing.Value = float64(existing.Quantity) * input.Price
			existing.UpdatedAt = time.Now()

			if err := db.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory"})
				return
			}

			var updated models.InventoryItem
			if err := db.Preload("Warehouse.Company").Preload("Responsible.Company").
				First(&updated, existing.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "preload after update failed", "details": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "inventory quantity updated", "item": updated})
			return

		}

		value := float64(input.Quantity) * input.Price
		newItem := models.InventoryItem{
			Name:          input.Name,
			Category:      input.Category,
			Condition:     input.Condition,
			Quantity:      input.Quantity,
			Price:         input.Price,
			Value:         value,
			WarehouseID:   input.WarehouseID,
			ResponsibleID: input.ResponsibleID,
			PurchaseDate:  purchaseDate,
		}

		if err := db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create inventory item"})
			return
		}

		var result models.InventoryItem
		if err := db.Preload("Warehouse.Company").Preload("Responsible.Company").First(&result, newItem.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inventory item with relations"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "inventory item created", "item": result})
	}
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.InventoryItem
		if err := db.Preload("Warehouse.Company").
			Preload("Responsible.Company").
			Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list inventory items"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func Detail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory item id"})
			return
		}

		var item models.InventoryItem
		if err := db.Preload("Warehouse.Company").
			Preload("Responsible.Company").
			First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "inventory item not found"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func Update(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
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
			Name          string   `json:"name"`
			Category      string   `json:"category"`
			Condition     string   `json:"condition"`
			Quantity      int      `json:"quantity"`
			WarehouseID   uint     `json:"warehouse_id"`
			ResponsibleID *uint    `json:"responsible_id"`
			PurchaseDate  *string  `json:"purchase_date"`
			Price         *float64 `json:"price"`
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
		if input.Price != nil {
			item.Price = *input.Price
		}
		if input.Price != nil || input.Quantity != 0 {
			item.Value = item.Price * float64(item.Quantity)
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

		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory item"})
			return
		}

		var result models.InventoryItem
		if err := db.Preload("Warehouse.Company").
			Preload("Responsible.Company").
			First(&result, item.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated inventory with relations"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
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
