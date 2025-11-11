package handlers

import (
	"context"
	"net/http"
	"task-manager/backend/internal/db"
	"task-manager/backend/internal/models"
	"task-manager/backend/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Регистрация нового пользователя
func Register(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Некорректные данные")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Ошибка шифрования пароля")
	}

	var id int
	err = db.Pool.QueryRow(context.Background(),
		"INSERT INTO users (username, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id",
		req.Username, string(hash), time.Now(),
	).Scan(&id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Имя пользователя уже занято")
	}

	token, err := utils.GenerateToken(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Ошибка генерации токена")
	}

	// Ставим cookie с JWT
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true на проде с HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusCreated, echo.Map{
		"username": req.Username,
	})
}

// Логин пользователя
func Login(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Некорректные данные")
	}

	var user models.User
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE username=$1", req.Username,
	).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Неверный логин или пароль")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.JSON(http.StatusUnauthorized, "Неверный логин или пароль")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Ошибка генерации токена")
	}

	// Ставим cookie с JWT
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true на проде с HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"username": req.Username,
	})
}

// (Опционально) Logout
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged out",
	})
}
