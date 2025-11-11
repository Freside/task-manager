package main

import (
	"log"
	"task-manager/backend/internal/db"
	"task-manager/backend/internal/handlers"
	customMiddleware "task-manager/backend/internal/middleware" //Тут переименовал чтобы было понятнее

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware" // Чтобы не было конфликта имён
)

func main() {
	// Подключение к БД
	db.Connect()
	defer db.Pool.Close()

	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// Публичный роут (для проверки состояния)
	e.GET("/api/health", handlers.Health)

	// Регистрация и логин
	e.POST("/api/register", handlers.Register)
	e.POST("/api/login", handlers.Login)

	// Приватные маршруты
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware)

	api.GET("/boards", handlers.GetBoards)
	api.POST("/boards", handlers.CreateBoard)

	api.GET("/boards/:board_id/columns", handlers.GetColumns)
	api.POST("/columns", handlers.CreateColumn)

	api.GET("/columns/:column_id/cards", handlers.GetCards)
	api.POST("/cards", handlers.CreateCard)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
