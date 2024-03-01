package main

import (
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := flow.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handle("/v1/healthcheck", http.HandlerFunc(app.healthcheckHandler), "GET")

	router.Handle("/v1/movies", http.HandlerFunc(app.listMoviesHandler), "GET")
	router.Handle("/v1/movies", http.HandlerFunc(app.createMovieHandler), "POST")

	router.Handle("/v1/movies/:id", http.HandlerFunc(app.showMovieHandler), "GET")
	router.Handle("/v1/movies/:id", http.HandlerFunc(app.updateMovieHandler), "PATCH")
	router.Handle("/v1/movies/:id", http.HandlerFunc(app.deleteMovieHandler), "DELETE")

	return app.recoverPanic(app.ratelimit(router))
}
