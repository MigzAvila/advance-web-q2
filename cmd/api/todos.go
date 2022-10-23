// Filename: cmd/api/todos.go

package main

import (
	"errors"
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

	// create a todo
	err = app.models.Todos.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create a Location header for the newly created resource/todos
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todos/%d", todo.ID))
	// write the json response with 201 - created status code with the body
	// being the todo data and the headers being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

// showTodoHandler for GET /v1/todos endpoints
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	//Utilize Utility Methods From helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific todo
	todo, err := app.models.Todos.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// write the data return by the Get method
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// updateTodoHandler for PUT /v1/todos/{id} endpoints
func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// get the id of the todo and update the todo
	// Utilize Utility Methods From helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fetch the original record from database
	todo, err := app.models.Todos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// create an input struct to hold the data read in from the client
	// Update input struct to use pointers because pointers have a default value of nil
	// if field remains nil then we know that the client is not interested in updating the field
	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"Description"`
		Completed   *bool   `json:"Completed"`
	}
	// Decode the data from the client
	err = app.readJSON(w, r, &input)

	// copy / update the fields / values in the todo variable using the fields in the input struct
	if err != nil {
		app.badResquestReponse(w, r, err)
		return
	}

	if input.Title != nil {
		todo.Title = *input.Title
	}

	if input.Description != nil {
		todo.Description = *input.Description
	}

	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	// validate the data provided by the client, if the validation fails,
	// then we send a 422 - Unprocessable responses to the client
	// Initialize a new validation error instance

	v := validator.New()

	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated todo record to the update method
	err = app.models.Todos.Update(todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// write the json response by Update
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

// deleteTodoHandler for DELETE /v1/todos/{id} endpoints
func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a delete of a specific todo
	// get the id of the todo and delete the todo
	// Utilize Utility Methods From helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// delete the todo from the database. send a 404 notFoundResponse status code to the client if there is no matching record
	// fetch the original record from database
	err = app.models.Todos.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//  return 200 status ok the client with a successful message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "todo successfully deleted"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

// listTodoHandler for GET /v1/todos endpoints (allows the client to see a listing of todos)
// based on a set of criteria
func (app *application) listTodosHandler(w http.ResponseWriter, r *http.Request) {
	// create an input struct to hold our query parameters
	var input struct {
		Title       string
		Description string
		Completed   bool
		data.Filters
	}

	// initialize a validator
	v := validator.New()

	// get the URL values in a map
	qs := r.URL.Query()

	// use the helper method to extract the values
	input.Title = app.readString(qs, "title", "")
	input.Description = app.readString(qs, "description", "")
	input.Completed = app.readBool(qs, "completed", false, v)

	// get the  page info
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)

	// get the sort info
	input.Filters.Sort = app.readString(qs, "sort", "id")

	// specific the allowed sort types
	input.Filters.SortList = []string{"id", "title", "description", "completed", "-id", "-title", "-description", "-completed"}

	// check for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	// Get a listing of all todos
	todos, metadata, err := app.models.Todos.GetAll(input.Title, input.Description, input.Completed, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todos": todos, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
