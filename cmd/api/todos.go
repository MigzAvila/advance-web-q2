// Filename: cmd/api/todos.go

package main

import (
	"fmt"
	"net/http"

	"todoapi.miguelavila.net/internals/data"
	"todoapi.miguelavila.net/internals/validator"
)

// createTodoHandler for POST v1/todos endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Completed   bool   `json:"completed"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badResquestReponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Todo struct
	todo := &data.Todo{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	}

	// Initialize a new instance of validator
	v := validator.New()

	// Check the errors maps if there were any errors validation
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// create a Location header for the newly created resource/school
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/schools/%d", todo.ID))
	// write the json response with 201 - created status code with the body
	// being the school data and the headers being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

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
