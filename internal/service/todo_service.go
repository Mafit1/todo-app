package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"todo-app/internal/dto"
	"todo-app/internal/models"
	"todo-app/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TodoService struct {
	repo  repository.TodoRepository
	redis *redis.Client
}

func NewTodoService(repo repository.TodoRepository, redisClient *redis.Client) *TodoService {
	return &TodoService{repo, redisClient}
}

func (s *TodoService) GetAll(ctx context.Context) ([]models.Todo, error) {
	key := "todos:all"
	ttl := time.Minute * 5

	return s.getListFromCacheOrDB(
		ctx,
		key,
		ttl,
		func() ([]models.Todo, error) {
			return s.repo.GetAll()
		},
	)
}

func (s *TodoService) Create(todo *models.Todo) error {
	return s.repo.Create(todo)
}

func (s *TodoService) GetById(ctx context.Context, id int) (*models.Todo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}

	key := fmt.Sprintf("todo:%d", id)
	ttl := time.Minute * 10

	return s.getFromCacheOrDB(
		ctx,
		key,
		ttl,
		func() (*models.Todo, error) {
			return s.repo.GetById(id)
		},
	)
}

func (s *TodoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete from DB: %w", err)
	}

	return s.invalidateTodoCache(ctx, id)
}

func (s *TodoService) Update(ctx context.Context, id int, req *dto.UpdateTodoRequest) (*models.Todo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}

	if req.Title != nil && len(*req.Title) > 100 {
		return nil, fmt.Errorf("title too long")
	}

	todo, err := s.repo.Update(id, req.Title, req.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to update record with id: %d in DB: %w", id, err)
	}

	err = s.invalidateTodoCache(ctx, id)
	if err != nil {
		fmt.Printf("warning: cache invalidation failed: %v", err)
	}

	return todo, nil
}

func (s *TodoService) getFromCacheOrDB(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetch func() (*models.Todo, error),
) (*models.Todo, error) {
	vals, err := s.redis.HGetAll(ctx, key).Result()
	if err == nil && len(vals) > 0 {
		todo := &models.Todo{}

		idStr := strings.TrimPrefix(key, "todo:")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			todo.ID = 0
		} else {
			todo.ID = id
		}

		todo.Title = vals["title"]
		completed, err := strconv.ParseBool(vals["completed"])
		if err != nil {
			completed = false
		}
		todo.Completed = completed

		fmt.Printf("Key: %s received from cache\n", key)
		return todo, nil
	}

	todo, err := fetch()
	if err != nil {
		return nil, err
	}

	_, err = s.redis.HSet(ctx, key, map[string]any{
		"title":     todo.Title,
		"completed": strconv.FormatBool(todo.Completed),
	}).Result()
	if err != nil {
		fmt.Printf("Ошибка записи в кэш Redis: %v\n", err)
	} else {
		s.redis.Expire(ctx, key, ttl)
	}

	fmt.Printf("Key: %s received from db\n", key)
	return todo, nil
}

func (s *TodoService) getListFromCacheOrDB(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetch func() ([]models.Todo, error),
) ([]models.Todo, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil && val != "" {
		var todos []models.Todo
		err = json.Unmarshal([]byte(val), &todos)
		if err == nil {
			fmt.Printf("Key: %s received from cache\n", key)
			return todos, nil
		}
	}

	todos, err := fetch()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(todos)
	if err == nil {
		s.redis.Set(ctx, key, data, ttl)
	}

	fmt.Printf("Key: %s received from db\n", key)
	return todos, nil
}

func (s *TodoService) invalidateTodoCache(ctx context.Context, id int) error {
	tx := s.redis.TxPipeline()
	tx.Del(ctx, fmt.Sprintf("todo:%d", id))
	tx.Del(ctx, "todos:all")
	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("redis transaction failed: %w", err)
	}
	return nil
}
