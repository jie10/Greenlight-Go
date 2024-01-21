package main

import (
	"errors"
	"github.com/alexedwards/flow"
	"net/http"
	"strconv"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := flow.Param(r.Context(), "id")

	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
