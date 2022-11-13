// Filename: cmd/api/routes
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Create a new HttpRouter router instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	
	router.HandlerFunc(http.MethodGet, "/v1/entries", app.listEntryHandler)
	
	router.HandlerFunc(http.MethodPost, "/v1/entries", app.createEntryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/entries/:id", app.showEntryHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/entries/:id", app.updateEntryHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/entries/:id", app.deleteEntryHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/stringrandom/:id", app.showRandomString)
	
	return app.recoverPanic(router)
}