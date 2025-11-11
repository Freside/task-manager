package models

import "time"

// Доска (Board)
type Board struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"` // новая строка
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Колонка (Column)
type Column struct {
	ID        int       `json:"id" db:"id"`
	BoardID   int       `json:"board_id" db:"board_id"`
	Title     string    `json:"title" db:"title"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Карточка (Card)
type Card struct {
	ID        int       `json:"id" db:"id"`
	ColumnID  int       `json:"column_id" db:"column_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
