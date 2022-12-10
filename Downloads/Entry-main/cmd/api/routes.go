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
	
	router.HandlerFunc(http.MethodGet, "/v1/entries", app.requirePermission("entries:read", app.listEntryHandler))
	router.HandlerFunc(http.MethodPost, "/v1/entries", app.requirePermission("entries:write", app.createEntryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/entries/:id", app.requirePermission("entries:read", app.showEntryHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/entries/:id", app.requirePermission("entries:write", app.updateEntryHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/entries/:id", app.requirePermission("entries:write", app.deleteEntryHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/stringrandom/:id", app.showRandomString)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}