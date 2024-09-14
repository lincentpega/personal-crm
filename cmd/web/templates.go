package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func (app *application) loadTemplates() error {
	app.templates = make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts := template.New(name)
		ts, err := ts.ParseFiles("./ui/html/base.html")
		if err != nil {
			return err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return err
		}

		app.templates[name] = ts
	}

	return nil
}

func (app *application) render(w http.ResponseWriter, name string, data interface{}) {
	ts, ok := app.templates[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}
