package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Генерация токена
func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // токен на 24 часа
		},
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET не задан в .env")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
