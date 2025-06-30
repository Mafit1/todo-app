package database

import (
	"context"
	"database/sql"
	"fmt"
)

func CreateTodoTable(ctx context.Context, db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS todo (
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            completed BOOLEAN DEFAULT FALSE
        );
    `
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create todo table: %w", err)
	}
	return nil
}
