package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nu12/audio-gonverter/internal/model"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.LoadSessionAndUser)

	mux.Get("/", app.Home)
	mux.Get("/status", app.Status)
	mux.Post("/upload", app.Upload)
	mux.Post("/convert", app.Convert)

	fileServer := http.FileServer(http.Dir(app.StaticFilesPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// SA1029: Users of WithValue should define their own types for keys.
type userID string

func (app *Config) LoadSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := app.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			log.Error(err)
			next.ServeHTTP(w, r)
			return
		}

		if session.IsNew {
			log.Debug("Creating new User for session")
			session.Values["user"] = model.NewUser().UUID
			if err := session.Save(r, w); err != nil {
				log.Error(err)
				next.ServeHTTP(w, r)
				return
			}
		}

		user := app.loadUser(session.Values["user"].(string))
		ctx := r.Context()
		newCtx := context.WithValue(ctx, userID("user"), user)
		sr := r.WithContext(newCtx)
		next.ServeHTTP(w, sr)
	})
}
