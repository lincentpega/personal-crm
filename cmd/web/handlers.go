package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, "home.html", nil)
}
