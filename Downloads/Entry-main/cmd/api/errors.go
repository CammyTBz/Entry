package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// We want to send JSON formatted error messages
func (app *application) errorResponse (w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	// Create the JSON response
	env := envelope{"error": message}
	err := app.writeJSON (w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Server error Response
func (app *application) serverErrorResponse (w http.ResponseWriter, r *http.Request, err error) {
	// log the error
	app.logError(r, err)
	// Prepare a message with the error
	message := "the server encountered a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Not Found Response
func (app *application) notFoundResponse (w http.ResponseWriter, r *http.Request) {
	// Create message
	message := "The requested resources could not be found."
	app.errorResponse (w, r, http.StatusNotFound, message)
}

// Not Allowed Response
func (app *application) methodNotAllowedResponse (w http.ResponseWriter, r *http.Request) {
	// Create message
	message := fmt.Sprintf("The %s method is not supported for this resource.", r.Method)
	app.errorResponse (w, r, http.StatusMethodNotAllowed, message)
}

// User Provided Bad Request
func (app *application) badRequestResponse (w http.ResponseWriter, r *http.Request, err error) {
	// Create message
	app.errorResponse (w, r, http.StatusBadRequest, err.Error())
}

// User provided an invalid Validation
func (app *application) failedValidationResponse (w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// EditError
func (app *application) editConflictResponse (w http.ResponseWriter, r *http.Request) {
	message := "unable to unable the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}