package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", handlers.Register)
		v1.POST("/login", handlers.Login)
		v1.POST("/verify-otp", handlers.VerifyOTP)
		v1.POST("/refresh", handlers.RefreshToken)
		v1.POST("/logout", handlers.Logout)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handlers.GetMe)
		}
	}
}
