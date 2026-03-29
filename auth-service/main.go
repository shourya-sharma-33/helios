package main

import (
    "github.com/gin-gonic/gin"
    "auth-service/config"
)

func main() {
    db := config.ConnectDB()
    _ = db // avoid unused error

    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Auth service running 🚀",
        })
    })

    r.Run(":8080")
}