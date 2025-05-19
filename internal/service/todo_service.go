package service

import (
	"fmt"
	"todo-app/internal/models"
	"todo-app/internal/repository"
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo}
}

func (s *TodoService) GetAll() ([]models.Todo, error) {
	return s.repo.GetAll()
}

func (s *TodoService) Create(todo *models.Todo) error {
	return s.repo.Create(todo)
}

func (s *TodoService) GetById(id int) (*models.Todo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.repo.GetById(id)
}
