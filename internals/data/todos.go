// Filename : internal/data/todos.go

package data

import (
	"database/sql"
)

// define a TodosModel object that wraps a sql.DB connection pool

type TodosModel struct {
	DB *sql.DB
}
