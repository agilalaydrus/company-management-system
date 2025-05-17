package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"metro-backend/internal/models"
	"os"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials (email not found)"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials (wrong password)"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		})

		jwtKey := []byte(os.Getenv("JWT_SECRET"))
		tokenString, _ := token.SignedString(jwtKey)

		c.JSON(200, LoginResponse{Token: tokenString})
	}
}
