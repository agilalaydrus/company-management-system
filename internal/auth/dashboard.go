package auth

import (
	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Metro Management System Dashboard",
	})
}
