// Filename: cmd/api/todos.go

package main

import (
	"fmt"
	"net/http"
)

// createTodoHandler for POST v1/todos endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("creating a new todo task")

}

// showTodoHandler for GET /v1/todos endpoints
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getting specific todo task")

}

// updateTodoHandler for PUT /v1/todos/{id} endpoints
func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("updating specific todo task for %v", id)
}

// deleteTodoHandler for DELETE /v1/todos/{id} endpoints
func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("deleting specific todo task for %v", id)

}

// listTodoHandler for GET /v1/todos endpoints (allows the client to see a listing of schools)
// based on a set of criteria
func (app *application) listTodosHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getting list of todo tasks")
}
