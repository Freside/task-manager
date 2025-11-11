package handlers

import (
	"context"
	"net/http"
	"task-manager/backend/internal/db"
	"task-manager/backend/internal/models"
	"time"

	"github.com/labstack/echo/v4"
)

// Получить все карточки по column_id
func GetCards(c echo.Context) error {
	columnID := c.Param("column_id")
	userID := c.Get("user_id").(int)

	rows, err := db.Pool.Query(context.Background(),
		`SELECT ca.id, ca.column_id, ca.title, ca.content, ca."order", ca.created_at
		 FROM cards ca
		 JOIN columns co ON co.id = ca.column_id
		 JOIN boards b ON b.id = co.board_id
		 WHERE ca.column_id=$1 AND b.user_id=$2
		 ORDER BY ca."order"`, columnID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(&card.ID, &card.ColumnID, &card.Title, &card.Content, &card.Order, &card.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		cards = append(cards, card)
	}

	return c.JSON(http.StatusOK, cards)
}

// Создать новую карточку
func CreateCard(c echo.Context) error {
	var card models.Card
	if err := c.Bind(&card); err != nil {
		return c.JSON(http.StatusBadRequest, "Некорректные данные")
	}

	card.CreatedAt = time.Now()

	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO cards (column_id, title, content, \"order\", created_at) VALUES ($1, $2, $3, $4, $5)",
		card.ColumnID, card.Title, card.Content, card.Order, card.CreatedAt,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, card)
}
