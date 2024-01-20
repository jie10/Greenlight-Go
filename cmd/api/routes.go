package main

import "github.com/alexedwards/flow"

func (app *application) routes() *flow.Mux {
	router := flow.New()

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler, "GET")

	return router
}
