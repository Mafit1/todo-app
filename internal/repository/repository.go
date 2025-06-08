package repository

import "todo-app/internal/models"

type TodoRepository interface {
	GetAll() ([]models.Todo, error)
	GetById(id int) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(id int, title *string, completed *bool) (*models.Todo, error)
	Delete(id int) error
}
