package main

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

type TemplateData struct {
	Commit string
}

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.page.gohtml").Delims("<<", ">>").ParseFiles(app.TemplatesPath + "index.page.gohtml")
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, TemplateData{Commit: app.Env["COMMIT"]})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	fileServer := http.FileServer(http.Dir(app.StaticFilesPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
