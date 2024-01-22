package main

import (
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() *flow.Mux {
	router := flow.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler, "GET")
	router.HandleFunc("/v1/movies", app.createMovieHandler, "POST")
	router.HandleFunc("/v1/movies/:id", app.showMovieHandler, "GET")

	return router
}
