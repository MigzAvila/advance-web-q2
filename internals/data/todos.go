// Filename : internal/data/todos.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todoapi.miguelavila.net/internals/validator"
)

type Todo struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
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
		INSERT INTO todos (title, description, completed)	
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
        SELECT id, title, description, completed, version
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

// Update() allows us to update a specific todo
// KEY: GO's http.server handles each request in its own goroutine
// Avoid data races
// A: Apples 3 buys 3 so 0 remains
// B: Apples 3 buys 2 so 1 remains
// USING Optimistic Locking to prevent multiple Optimistic sql
func (m TodosModel) Update(todo *Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3, version = version + 1
		WHERE id = $4
		AND version = $5
		RETURNING version
	`
	// Create a context
	// Time starts when the context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// cleanup the context to prevent memory leaks
	defer cancel()
	args := []interface{}{
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.ID,
		todo.Version,
	}

	// check for edit conflict
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete() allows us to delete a specific Todo
func (m TodosModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return nil
	}

	// Create the query for deleting a specific todo
	query := `
			DELETE FROM todos
			WHERE id = $1
		`

	// Create a context
	// Time starts when the context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// cleanup the context to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		// Check error type
		return err
	}

	// Check how many records were deleted by the query
	rows, err := result.RowsAffected()

	if err != nil {
		// Check error type
		return err
	}

	// check if no records were deleted
	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil

}

// func GetAll() method returns a list of all todo sorted by id
func (m TodosModel) GetAll(title string, description string, completed bool, filters Filters) ([]*Todo, Metadata, error) {
	// construct the query
	query := fmt.Sprintf(`
		 SELECT
		 		COUNT(*) OVER(),
				id, title, description, completed, version
				FROM todos
				WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
				AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
				AND ((completed = $3) OR $3 = false)
				ORDER BY %s %s, id ASC
				LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())

	// query := fmt.Sprintf(`
	// 		SELECT
	//         COUNT(*) OVER(),
	//         id, title, description, completed, version
	// 		FROM todos
	// 		WHERE (LOWER(title) = LOWER($1) or $1 = '')
	// 		AND (LOWER(description) = LOWER($2) or $2 = '')
	// 		AND ((completed = $3) OR $3 = false)
	// 		ORDER BY %s %s, id ASC
	// 		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())

	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// cleanup the context to prevent memory leaks
	defer cancel()

	args := []interface{}{title, description, completed, filters.limit(), filters.offset()}

	// execute the query
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		// Check error type
		return nil, Metadata{}, err
	}

	// cleanup the rows to prevent memory leaks
	defer rows.Close()

	totalRecords := 0

	// initialize an empty slice
	todos := []*Todo{}

	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// add the todo to the slice
		todos = append(todos, &todo)
	}
	// check for errors after looping the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculatesMetadata(totalRecords, filters.Page, filters.PageSize)

	// return the slice of Todos
	return todos, metadata, nil
}
