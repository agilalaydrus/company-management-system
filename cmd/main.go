package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"metro-backend/internal/db"
	"metro-backend/routes"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	db.AutoMigrate(database)

	router := gin.Default()

	routes.RegisterRoutes(router, database)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("✅ Server running on port :" + port)
	router.Run("0.0.0.0:" + port)

}
