package controllers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ================= REGISTER =================
func Register(c *gin.Context) {

	var input struct {
		Name     string `json:"name" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}

	// rate limit
	cooldown := "cooldown:" + input.Email
	if exists, _ := config.Rdb.Exists(config.Ctx, cooldown).Result(); exists == 1 {
		c.JSON(429, gin.H{"success": false, "error": "Too many requests"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	token := utils.GenerateToken()

	data := map[string]string{
		"name":     input.Name,
		"email":    input.Email,
		"password": string(hashed),
	}

	jsonData, _ := json.Marshal(data)

	config.Rdb.Set(config.Ctx, "verify:"+token, jsonData, time.Minute*10)
	config.Rdb.Set(config.Ctx, cooldown, "1", time.Minute)

	utils.SendEmail(input.Email, token)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Verification email sent",
	})
}

// ================= VERIFY =================
func Verify(c *gin.Context) {

	token := c.Param("token")

	data, err := config.Rdb.Get(config.Ctx, "verify:"+token).Result()
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid/expired token"})
		return
	}

	var user map[string]string
	json.Unmarshal([]byte(data), &user)

	_, err = config.DB.Exec(`
		INSERT INTO users (email, username, name, password)
		VALUES ($1, $2, $3, $4)
	`,
		user["email"],
		user["name"],
		user["name"],
		user["password"],
	)

	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": "DB error: " + err.Error()})
		return
	}

	config.Rdb.Del(config.Ctx, "verify:"+token)

	c.JSON(200, gin.H{
		"success": true,
		"message": "User verified",
	})
}

// ================= LOGIN (SEND OTP) =================
func Login(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔥 RATE LIMIT (email + ip)
	ip := c.ClientIP()
	rateKey := "login_rate:" + input.Email + ":" + ip

	exists, _ := config.Rdb.Exists(config.Ctx, rateKey).Result()
	if exists == 1 {
		c.JSON(429, gin.H{"error": "Too many requests, try later"})
		return
	}

	// 🔍 FIND USER
	var user models.User
	err := config.DB.Get(&user, `
		SELECT id, email, password, name
		FROM users
		WHERE email=$1
	`, input.Email)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid credentials"})
		return
	}

	// 🔐 PASSWORD CHECK
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid credentials"})
		return
	}

	// 🔢 GENERATE OTP (6 digit)
	otp := fmt.Sprintf("%06d", rand.Intn(900000)+100000)

	otpKey := "otp:" + input.Email

	// store in redis (5 min)
	config.Rdb.Set(config.Ctx, otpKey, otp, 5*time.Minute)

	// set rate limit (1 min)
	config.Rdb.Set(config.Ctx, rateKey, "1", time.Minute)

	// 📧 SEND EMAIL (or log)
	fmt.Printf("OTP FOR %s: %s\n", input.Email, otp)
	utils.SendOTP(input.Email, otp)

	c.JSON(200, gin.H{
		"message": "If email is valid, OTP sent (valid 5 min)",
	})
}

// ================= VERIFY OTP + LOGIN COMPLETE =================
func VerifyOTP(c *gin.Context) {

	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if input.Email == "" || input.OTP == "" {
		c.JSON(400, gin.H{"error": "Provide all fields"})
		return
	}

	// 🔍 GET OTP FROM REDIS
	otpKey := "otp:" + input.Email

	storedOTP, err := config.Rdb.Get(config.Ctx, otpKey).Result()
	if err != nil {
		c.JSON(400, gin.H{"error": "OTP expired"})
		return
	}

	// ❌ WRONG OTP
	if storedOTP != input.OTP {
		c.JSON(400, gin.H{"error": "Invalid OTP"})
		return
	}

	// ✅ DELETE OTP
	config.Rdb.Del(config.Ctx, otpKey)

	// 🔍 FIND USER AGAIN
	var user models.User
	err = config.DB.Get(&user, `
		SELECT id, email, name
		FROM users
		WHERE email=$1
	`, input.Email)

	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// 🔐 GENERATE TOKENS
	accessToken, _ := utils.GenerateJWT(user.ID, time.Minute*15)
	refreshToken, _ := utils.GenerateJWT(user.ID, time.Hour*24*7)

	// 🧠 STORE REFRESH TOKEN IN REDIS (SESSION)
	refreshKey := "refresh:" + user.ID
	config.Rdb.Set(config.Ctx, refreshKey, refreshToken, 7*24*time.Hour)

	c.JSON(200, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
