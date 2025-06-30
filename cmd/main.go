package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todo-app/config"
	"todo-app/internal/database"
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

	c := config.LoadConfig()

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := database.CreateTodoTable(ctx, db); err != nil {
		log.Fatalf("failed to create todo table: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			c.RedisHost,
			c.RedisPort,
		),
		Password: c.RedisPassword,
		DB:       0,
	})

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

	e.Logger.Fatal(e.Start(c.ServerPort))
}
