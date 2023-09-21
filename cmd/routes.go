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
	mux.Get("/delete", app.Delete)
	mux.Get("/delete-all", app.DeleteAll)

	fileServer := http.FileServer(http.Dir(app.StaticFilesPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// TODO: configure
	convertedServer := http.FileServer(http.Dir("/tmp/"))
	mux.Handle("/download/*", http.StripPrefix("/download", convertedServer))

	return mux
}

// SA1029: Users of WithValue should define their own types for keys.
type userID string

func (app *Config) CreateSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := app.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			log.Error(err)
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
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
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		if err := app.saveUser(&user); err != nil {
			log.Error(err)
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
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
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}

		user, err := app.loadUser(session.Values["user"].(string))
		if err != nil {
			log.Error(err)
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		newCtx := context.WithValue(ctx, userID("user"), user)
		sr := r.WithContext(newCtx)
		next.ServeHTTP(w, sr)
	})
}

func (app *Config) write(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
