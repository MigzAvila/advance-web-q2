// Filename : internal/data/todos.go

package data

import (
	"database/sql"
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
