package routes

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	r.POST("/register", controllers.Register)
	r.GET("/verify/:token", controllers.Verify)
	r.POST("/login", controllers.Login)

	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Protected route"})
	})
}
