// Filename/cmd/api/entry.go

package main

import (
	"fmt"
	"net/http"
	"errors"

	"kriol.camerontillett.net/internal/data"
	"kriol.camerontillett.net/internal/validator"
)

//createEntryHandler for the "POST /v1/entry" endpoint
//Change 
func (app *application) createEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Our Target Decode destination
	var input struct {
		Name string `json:"name"`
		Level string `json:"level"`
		Contact string `json:"contact"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Website string `json:"website"`
		Address string `json:"address"`
		Mode []string `json:"mode"`
	}
	// Initialize a new json.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Entry struct.
	entries := &data.Entry{
		Name: input.Name,
		Level: input.Level,
		Contact: input.Contact,
		Phone: input.Phone,
		Email: input.Email,
		Website: input.Website,
		Address: input.Address,
		Mode: input.Mode,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// check the map to see if there were validation errors
	if data.ValidateEntries(v, entries); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Create an entry
	err = app.models.Entry.Insert(entries)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Create a Location header for the newly created resource/school
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/entries/%d", entries.ID))
	// Write the JSON response with 201 - Created status code with the body
	// being the entry data and the header being the header map
	err = app.writeJSON(w, http.StatusCreated, envelope{"entries":entries}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//createEntryHandler for the "GET /v1/entry/:id" endpoint
func (app *application) showEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the value of the "id" parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific entry
	entries, err := app.models.Entry.Get(id)
	// Handle Errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return 
	}

	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"entries": entries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateEntryHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the school that need updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the original record from the database
	entries, err := app.models.Entry.Get(id)
	// Handle Errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return 
	}
	// Create an input struct to hold data read from the Client
	// Update the input struct to use pointers for default value of nil
	// If field remains nil we know it wasnt updated
	var input struct {
		Name 	*string  `json:"name"`
		Level 	*string  `json:"level"`
		Contact *string  `json:"contact"`
		Phone 	*string  `json:"phone"`
		Email 	*string  `json:"email"`
		Website *string  `json:"website"`
		Address *string  `json:"address"`
		Mode 	[]string `json:"mode"`
	}
	
	// Initialize a new json.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check for updates
	if input.Name != nil {
		entries.Name = *input.Name
	}
	if input.Level != nil {
		entries.Level = *input.Level
	}
	if input.Contact != nil {
		entries.Contact = *input.Contact
	}
	if input.Phone != nil {
		entries.Phone = *input.Phone
	}
	if input.Email != nil {
		entries.Email = *input.Email
	}
	if input.Website != nil {
		entries.Website = *input.Website
	}
	if input.Address != nil {
		entries.Address = *input.Address
	}
	if input.Mode != nil {
		entries.Mode = input.Mode
	}

	// Perform validation on the updated entry. If fails the we send a 422 - unprocessable response
	// Initialize validator instance
	v := validator.New()

	// check the map to see if there were validation errors
	if data.ValidateEntries(v, entries); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	
	// Pass the Updated Entry record to the Update () method
	err = app.models.Entry.Update(entries)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"entries": entries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the entry that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the School from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Entry.Delete(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "entry successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Create an input struct to hold our query parameter
	var input struct {
		Name string
		Level string
		Mode []string
		data.Filters
	}

	// initialize a validator
	v := validator.New()
	// Get the URL values map
	qs := r.URL.Query()
	// Use the helper methods to extract the values
	input.Name = app.readString(qs, "name", "")
	input.Level = app.readString(qs, "level", "")
	input.Mode = app.readCSV(qs, "mode", []string{})
	// Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specific the allowed sort values
	input.Filters.SortList = []string{"id", "name", "level", "-id", "-name", "-level"}
	// Check for validation errors
	if data.ValidateFilters(v, input.Filters);!v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get a listing of all the entries
	entries, metadata, err := app.models.Entry.GetAll(input.Name, input.Level, input.Mode, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSN response contain all the entries
	err = app.writeJSON(w, http.StatusOK, envelope{"entries": entries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// func (app *application) showRandomString (w http.ResponseWriter, r *http.Request) {

// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	integer := int(id)
// 	tools := &data.Tools{}

// 	random := tools.GenerateRandomString(integer)
// 	data := envelope{
// 		"Here is your randomize string": random,
// 		"Your :id is ":                   integer,
// 	}
// 	err = app.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// }