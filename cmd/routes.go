package main

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/user"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.CreateSessionAndUser)
	mux.Use(app.LoadSessionAndUser)
	mux.Use(app.StatusCheck)

	mux.Get("/", app.Home)
	mux.Post("/upload", app.Upload)
	mux.Post("/convert", app.Convert)
	mux.Get("/delete", app.Delete)
	mux.Get("/delete-all", app.DeleteAll)
	mux.Get("/download", app.Download)

	fileServer := http.FileServer(http.Dir(app.StaticFilesPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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

		user := user.New()
		session.Values["user"] = user.UUID
		log.Debug("Created new User for session: " + user.UUID)

		session.Options.MaxAge = 3600 // 1 hour
		if err := session.Save(r, w); err != nil {
			log.Error(err)
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		if err := app.saveUser(user); err != nil {
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
		e, err := app.DatabaseRepo.Exist(session.Values["user"].(string))
		if !e && err == nil {
			u := &user.User{UUID: session.Values["user"].(string), Files: []*file.File{}}
			err2 := app.saveUser(u)
			if err2 != nil {
				log.Error(err)
				app.write(w, err2.Error(), http.StatusInternalServerError)
				return
			}

		} else if err != nil {
			log.Error(err)
			app.write(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := app.loadUser(session.Values["user"].(string))
		if err != nil {
			log.Error(err)
			app.write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		sr := r.WithContext(user.ToContext(r.Context()))
		next.ServeHTTP(w, sr)
	})
}

func (app *Config) StatusCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := user.FromRequest(r)

		if user.IsConverting {
			app.render(w, "status.page.gohtml", TemplateData{Messages: []string{"Converting"}, Commit: app.Env["COMMIT"]})
			return
		}
		if user.IsUploading {
			app.render(w, "status.page.gohtml", TemplateData{Messages: []string{"Uploading"}, Commit: app.Env["COMMIT"]})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Config) write(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Config) render(w http.ResponseWriter, templateName string, td TemplateData) {
	t, err := template.New(templateName).ParseFiles(app.TemplatesPath+"base.layout.gohtml", app.TemplatesPath+templateName)
	if err != nil {
		log.Error(err)
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, td); err != nil {
		log.Error(err)
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
