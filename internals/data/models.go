// Filename: internal/data/models.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// A wrapper for out data models
type Models struct {
	Todos TodosModel
}

// NewModels() allows us to create new models
func NewModels(db *sql.DB) *Models {
	return &Models{
		Todos: TodosModel{DB: db},
	}
}

// insert() allows us to create a new Todo
func (m TodosModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (title, description, complete)	
		VALUES ($1, $2, $3)
		RETURNING id, create_at, version
	`
	// Create a context
	// Time starts when the context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// cleanup the context to prevent memory leaks
	defer cancel()

	// collect data fields into a slice
	args := []interface{}{
		todo.Title,
		todo.Description,
		todo.Completed,
	}
	// run query ... -> expand the slice
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}
