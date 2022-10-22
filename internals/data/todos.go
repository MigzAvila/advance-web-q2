// Filename : internal/data/todos.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"todoapi.miguelavila.net/internals/validator"
)

type Todo struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	Version     int32     `json:"version"`
}

// define a TodosModel object that wraps a sql.DB connection pool
type TodosModel struct {
	DB *sql.DB
}

func ValidateTodo(v *validator.Validator, Todo *Todo) {

	v.Check(Todo.Title != "", "title", "must be provided")
	v.Check(len(Todo.Title) <= 100, "title", "must be no more than 100 characters")

	v.Check(Todo.Description != "", "description", "must be provided")
	v.Check(len(Todo.Description) <= 1000, "description", "must be no more than 1000 characters")

	v.Check(Todo.Completed || !Todo.Completed, "completed", "must be a bool")

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

// Get() allows us to retrieve a specific todo
func (m TodosModel) Get(id int64) (*Todo, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query for getting a specific todo
	query := `
        SELECT id, title, description, complete, version
        FROM todos
        WHERE id = $1
    `
	// declare a todo variable and run query
	var todo Todo
	// Create a context
	// Time starts when the context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// cleanup the context to prevent memory leaks
	defer cancel()

	// Execute the query
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.Version,
	)

	if err != nil {
		// Check error type
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	// Success
	return &todo, nil
}
