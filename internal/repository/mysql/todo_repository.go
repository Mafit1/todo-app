package mysql

import (
	"database/sql"
	"fmt"
	"strings"
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

func (r *MySQLTodoRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM todo WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	return nil
}

func (r *MySQLTodoRepository) Update(id int, title *string, completed *bool) (*models.Todo, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := "UPDATE todo SET "
	args := []any{}
	updates := []string{}

	if title != nil {
		updates = append(updates, "title = ?")
		args = append(args, title)
	}
	if completed != nil {
		updates = append(updates, "completed = ?")
		args = append(args, completed)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo with id %d not found", id)
	}

	row := tx.QueryRow("SELECT id, title, completed FROM todo WHERE id = ?", id)
	var todo models.Todo

	if err := row.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo, nil
}
