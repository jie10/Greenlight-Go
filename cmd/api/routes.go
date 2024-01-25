package main

import (
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() *flow.Mux {
	router := flow.New()

	router.NotFound = app.recoverPanic(http.HandlerFunc(app.notFoundResponse))
	router.MethodNotAllowed = app.recoverPanic(http.HandlerFunc(app.methodNotAllowedResponse))

	router.Handle("/v1/healthcheck", app.recoverPanic(http.HandlerFunc(app.healthcheckHandler)), "GET")
	router.Handle("/v1/movies", app.recoverPanic(http.HandlerFunc(app.createMovieHandler)), "POST")
	router.Handle("/v1/movies/:id", app.recoverPanic(http.HandlerFunc(app.showMovieHandler)), "GET")

	return router
}
