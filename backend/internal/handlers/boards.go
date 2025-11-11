package handlers

import (
	"context"
	"net/http"
	"task-manager/backend/internal/db"
	"task-manager/backend/internal/models"
	"time"

	"github.com/labstack/echo/v4"
)

// Получить все доски
func GetBoards(c echo.Context) error {
	userID := c.Get("user_id").(int) // берём из JWT

	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, user_id, title, created_at FROM boards WHERE user_id=$1 ORDER BY id DESC", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.UserID, &b.Title, &b.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		boards = append(boards, b)
	}

	return c.JSON(http.StatusOK, boards)
}

// Создать новую доску
func CreateBoard(c echo.Context) error {
	userID := c.Get("user_id").(int) // берём из JWT

	var board models.Board
	if err := c.Bind(&board); err != nil {
		return c.JSON(http.StatusBadRequest, "Некорректные данные")
	}

	board.UserID = userID
	board.CreatedAt = time.Now()

	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO boards (user_id, title, created_at) VALUES ($1, $2, $3) RETURNING id",
		board.UserID, board.Title, board.CreatedAt,
	).Scan(&board.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, board)
}
