package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"todo-app/internal/handler"
	"todo-app/internal/repository/mysql"
	"todo-app/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Загрузка .env (если есть)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Подключение к MySQL
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация слоёв
	todoRepo := mysql.NewMySQLTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	// Настройка Echo
	e := echo.New()

	// Роуты
	e.GET("/todo", todoHandler.GetAll)
	e.POST("/todo", todoHandler.Create)
	e.GET("/todo/:id", todoHandler.GetByID)
	// e.PUT("/todo/:id", todoHandler.Update)
	// e.DELETE("/todo/:id", todoHandler.Delete)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
