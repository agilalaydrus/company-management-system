package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"metro-backend/internal/models"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to hash password"})
			return
		}

		user := models.User{
			Email:    req.Email,
			Password: string(hashed),
			Role:     req.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(201, gin.H{"message": "user registered successfully"})
	}
}
