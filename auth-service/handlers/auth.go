package handlers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================= REGISTER =================
func Register(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or DB error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// ================= LOGIN (SEND OTP) =================
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPassword(body.Password, user.Password) {
		c.JSON(400, gin.H{"error": "Invalid credentials"})
		return
	}

	// 🔥 Generate OTP (like frontend flow)
	otp := utils.GenerateOTP()
	user.OTP = otp
	config.DB.Save(&user)

	// TODO: send email (for now print)
	fmt.Println("OTP:", otp)

	c.JSON(200, gin.H{
		"message": "OTP sent to email",
	})
}

// ================= VERIFY OTP (MOST IMPORTANT) =================
func VerifyOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != body.OTP {
		c.JSON(400, gin.H{"error": "Invalid OTP"})
		return
	}

	user.IsVerified = true
	user.OTP = ""
	config.DB.Save(&user)

	// 🔥 Generate tokens
	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	// 🔥 Set cookies (IMPORTANT)
	c.SetCookie("access_token", accessToken, 900, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "OTP verified",
	})
}

// ================= GET USER (/me) =================
func GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	config.DB.First(&user, userID)

	c.JSON(200, gin.H{
		"user": user,
	})
}

// ================= REFRESH TOKEN (AUTO LOGIN) =================
func RefreshToken(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "No refresh token"})
		return
	}

	userIDStr, err := utils.ValidateRefreshToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	var userID uint
	fmt.Sscanf(userIDStr, "%d", &userID)

	newAccess, _ := utils.GenerateAccessToken(userID)

	c.SetCookie("access_token", newAccess, 900, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "Token refreshed",
	})
}

// ================= LOGOUT =================
func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "Logged out",
	})
}
