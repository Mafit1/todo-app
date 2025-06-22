package main

import (
	"context"
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
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_ADDRESS"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}

	todoRepo := mysql.NewMySQLTodoRepository(db)
	todoService := service.NewTodoService(todoRepo, redisClient)
	todoHandler := handler.NewTodoHandler(todoService)

	e := echo.New()

	e.GET("/todo", todoHandler.GetAll)
	e.POST("/todo", todoHandler.Create)
	e.GET("/todo/:id", todoHandler.GetByID)
	e.PUT("/todo/:id", todoHandler.Update)
	e.PATCH("/todo/:id", todoHandler.Update)
	e.DELETE("/todo/:id", todoHandler.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}
