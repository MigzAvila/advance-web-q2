// Filename : cmd/api/errors.go

package main

import (
	"net/http"
)

// Log errors
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// Send JSON-formatted error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	// create the json response
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// Server error message
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log the error
	app.logError(r, err)
	//prepare a message with error
	message := "the server encountered an problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// User passed a bad request
func (app *application) badResquestReponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Edit Conflict validation errors
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}