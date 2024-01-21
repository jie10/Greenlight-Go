package main

import (
	"fmt"
	"github.com/alexedwards/flow"
	"net/http"
	"strconv"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := flow.Param(r.Context(), "id")

	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, err = fmt.Fprintf(w, "show the details of movie %d\n", id)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
}
