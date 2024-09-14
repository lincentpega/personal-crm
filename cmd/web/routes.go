package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) route() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))

	return dynamic.Then(mux)
}
