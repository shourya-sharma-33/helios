package controllers

import (
	"auth-service/config"
	"auth-service/utils"
	"encoding/json"
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

// ================= LOGIN =================
func Login(c *gin.Context) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&input)

	var hashed string
	var userID string

	err := config.DB.QueryRow(`
		SELECT id, password
		FROM users
		WHERE email=$1
	`, input.Email).Scan(&userID, &hashed)

	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input.Password)) != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateJWT(userID)

	c.JSON(200, gin.H{
		"success": true,
		"token":   token,
	})
}
