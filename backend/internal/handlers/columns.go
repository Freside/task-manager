package handlers

import (
	"context"
	"net/http"
	"task-manager/backend/internal/db"
	"task-manager/backend/internal/models"
	"time"

	"github.com/labstack/echo/v4"
)

// Получить все колонки по board_id
func GetColumns(c echo.Context) error {
	boardID := c.Param("board_id")
	userID := c.Get("user_id").(int)

	rows, err := db.Pool.Query(context.Background(),
		`SELECT c.id, c.board_id, c.title, c."order", c.created_at
		 FROM columns c
		 JOIN boards b ON b.id = c.board_id
		 WHERE c.board_id=$1 AND b.user_id=$2
		 ORDER BY c."order"`, boardID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var columns []models.Column
	for rows.Next() {
		var col models.Column
		if err := rows.Scan(&col.ID, &col.BoardID, &col.Title, &col.Order, &col.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		columns = append(columns, col)
	}

	return c.JSON(http.StatusOK, columns)
}

// Создать новую колонку
func CreateColumn(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var col models.Column
	if err := c.Bind(&col); err != nil {
		return c.JSON(http.StatusBadRequest, "Некорректные данные")
	}

	// Проверяем, что доска принадлежит текущему пользователю
	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM boards WHERE id=$1 AND user_id=$2)", col.BoardID, userID,
	).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, "Доска не найдена или недоступна")
	}

	col.CreatedAt = time.Now()
	err = db.Pool.QueryRow(context.Background(),
		"INSERT INTO columns (board_id, title, \"order\", created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		col.BoardID, col.Title, col.Order, col.CreatedAt,
	).Scan(&col.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, col)
}
