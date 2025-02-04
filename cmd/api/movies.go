package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.artc0des.com/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Nosferatu",
		Year:      2025,
		Runtime:   120,
		Genres:    []string{"Blood", "Chickens", "Rats"},
		Version:   1,
	}

	var data = envelope{"movie": movie}

	err = app.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "unable to process your request", http.StatusInternalServerError)
		return
	}
}
