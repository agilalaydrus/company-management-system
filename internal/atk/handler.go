package atk

import (
	"net/http"
	"strconv"

	"metro-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name          string `json:"name" binding:"required"`
			Unit          string `json:"unit"`
			Quantity      int    `json:"quantity" binding:"required"`
			WarehouseID   uint   `json:"warehouse_id" binding:"required"`
			ResponsibleID *uint  `json:"responsible_id"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		atk := models.ATKItem{
			Name:          input.Name,
			Unit:          input.Unit,
			Quantity:      input.Quantity,
			WarehouseID:   input.WarehouseID,
			ResponsibleID: input.ResponsibleID,
		}

		if err := db.Create(&atk).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ATK item"})
			return
		}

		c.JSON(http.StatusCreated, atk)
	}
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var atks []models.ATKItem
		if err := db.Preload("Warehouse").Preload("Responsible").Find(&atks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ATK items"})
			return
		}
		c.JSON(http.StatusOK, atks)
	}
}

func Detail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ATK item id"})
			return
		}

		var atk models.ATKItem
		if err := db.Preload("Warehouse").Preload("Responsible").First(&atk, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ATK item not found"})
			return
		}

		c.JSON(http.StatusOK, atk)
	}
}

func Update(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ATK item id"})
			return
		}

		var atk models.ATKItem
		if err := db.First(&atk, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ATK item not found"})
			return
		}

		var input struct {
			Name          string `json:"name"`
			Unit          string `json:"unit"`
			Quantity      int    `json:"quantity"`
			WarehouseID   uint   `json:"warehouse_id"`
			ResponsibleID *uint  `json:"responsible_id"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name != "" {
			atk.Name = input.Name
		}
		if input.Unit != "" {
			atk.Unit = input.Unit
		}
		if input.Quantity != 0 {
			atk.Quantity = input.Quantity
		}
		if input.WarehouseID != 0 {
			atk.WarehouseID = input.WarehouseID
		}
		atk.ResponsibleID = input.ResponsibleID

		if err := db.Save(&atk).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ATK item"})
			return
		}
		c.JSON(http.StatusOK, atk)
	}
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ATK item id"})
			return
		}

		if err := db.Delete(&models.ATKItem{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ATK item"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
