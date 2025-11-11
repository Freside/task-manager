package middleware

import (
	"net/http"
	"os"
	"task-manager/backend/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Нет токена")
		}

		tokenStr := cookie.Value
		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Неверный токен")
		}

		//сохраняем userid в контекст
		c.Set("user_id", claims.UserID)

		return next(c)
	}
}
