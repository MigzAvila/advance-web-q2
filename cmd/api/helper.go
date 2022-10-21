// Filename: cmd/api/helpers

package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Convert our map into a JSON object
	js, err := json.Marshal(data)
	// Format the JSON object for cmd -- Takes more resources than printing it normally
	// js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	// Add a newline to make viewing on the terminal easier
	js = append(js, '\n')

	// Add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Specify that we will serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Write the I I byte slice containing the JSON response body
	w.Write(js)
	return nil
}
