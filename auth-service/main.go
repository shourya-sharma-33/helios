package main

import (
	"auth-service/config"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	config.ConnectRedis()

	r := gin.Default()

	routes.AuthRoutes(r)

	r.Run(":8081")
}
