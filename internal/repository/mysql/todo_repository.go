package mysql

import (
	"database/sql"
	"fmt"
	"todo-app/internal/models"
)

type MySQLTodoRepository struct {
	db *sql.DB
}

func NewMySQLTodoRepository(db *sql.DB) *MySQLTodoRepository {
	return &MySQLTodoRepository{db}
}

func (r *MySQLTodoRepository) GetAll() ([]models.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, completed FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *MySQLTodoRepository) Create(todo *models.Todo) error {
	result, err := r.db.Exec(
		"INSERT INTO todo (title, completed) VALUES (?, ?)",
		todo.Title, todo.Completed,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	todo.ID = int(id)
	return nil
}

// Delete implements repository.TodoRepository.
func (r *MySQLTodoRepository) Delete(id int) error {
	panic("unimplemented")
}

func (r *MySQLTodoRepository) GetById(id int) (*models.Todo, error) {
	row := r.db.QueryRow("SELECT id, title, completed FROM todo WHERE id = ?", id)

	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.Completed)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("todo not found")
	} else if err != nil {
		return nil, err
	}
	return &todo, nil
}

// Update implements repository.TodoRepository.
func (r *MySQLTodoRepository) Update(todo *models.Todo) error {
	panic("unimplemented")
}
