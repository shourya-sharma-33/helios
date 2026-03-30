package main

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	config.ConnectDB()

	// Auto Migration
	config.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	// Setup Routes
	routes.AuthRoutes(r)

	r.Run(":8081")
}
