package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("supersecret")

func GenerateAccessToken(userID uint) (string, error) {
	return generateToken(userID, time.Minute*15)
}

func GenerateRefreshToken(userID uint) (string, error) {
	return generateToken(userID, time.Hour*24*7)
}

func generateToken(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userID),
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

func ValidateToken(tokenStr string) (string, error) {
	return validate(tokenStr)
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	return validate(tokenStr)
}

func validate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SECRET, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", errors.New("invalid token")
}
