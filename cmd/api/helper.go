// Filename: cmd/api/helpers

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Convert our map into a JSON object
	// js, err := json.Marshal(data)
	// Format the JSON object for cmd -- Takes more resources than printing it normally
	js, err := json.MarshalIndent(data, "", "\t")

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

// Utility function for reading ID in Endpoint
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Use the param
	// Use the ParamsFormContext
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}

	return id, nil

}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// use http.MaxBytesReader() to limit size of response body
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	//Decode the response body into the target destination
	err := dec.Decode(dst)

	// Check for bad responses
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		// Switch to check for errors
		switch {
		// Check for syntaxError
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON body (at character %d)", syntaxError.Offset)
		// Check for wrong body passed by client
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON body")
		// Check for wrong types passed by client
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %q)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Check for unmappable fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Body size to large
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not exceed %d bytes", maxBytes)
		// Pass non-nil error
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Call Decode() again
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
