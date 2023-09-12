package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nu12/audio-gonverter/internal/model"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.CreateSessionAndUser)
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

func (app *Config) CreateSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := app.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Session error. Cleaning the cache may solve the issue"))
			return
		}

		if !session.IsNew {
			next.ServeHTTP(w, r)
			return
		}

		user := model.NewUser()
		session.Values["user"] = user.UUID
		log.Debug("Created new User for session: " + user.UUID)

		session.Options.MaxAge = 3600 // 1 hour
		if err := session.Save(r, w); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Session error. Cleaning the cache may solve the issue"))
			return
		}
		if err := app.saveUser(&user); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Session error. Cleaning the cache may solve the issue"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Config) LoadSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := app.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Session error. Cleaning the cache may solve the issue"))
			return
		}

		user, err := app.loadUser(session.Values["user"].(string))
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Session error. Cleaning the cache may solve the issue"))
			return
		}
		ctx := r.Context()
		newCtx := context.WithValue(ctx, userID("user"), user)
		sr := r.WithContext(newCtx)
		next.ServeHTTP(w, sr)
	})
}
