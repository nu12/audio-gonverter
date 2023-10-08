package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/helper"
	"github.com/nu12/audio-gonverter/internal/user"
)

type Middleware struct {
	Config *config.Config
}

func Routes(app *config.Config) http.Handler {
	mux := chi.NewRouter()

	m := Middleware{Config: app}
	mux.Use(m.CreateSessionAndUser)
	mux.Use(m.LoadSessionAndUser)
	mux.Use(m.StatusCheck)

	handler := Handler{Config: app}
	mux.Get("/", handler.Home)
	mux.Post("/upload", handler.Upload)
	mux.Post("/convert", handler.Convert)
	mux.Get("/delete", handler.Delete)
	mux.Get("/delete-all", handler.DeleteAll)
	mux.Get("/download", handler.Download)

	fileServer := http.FileServer(http.Dir(app.StaticFilesPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func (m *Middleware) CreateSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := helper.WithConfig(m.Config)

		session, err := m.Config.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			m.Config.Log.Error(err)
			write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}

		if !session.IsNew {
			next.ServeHTTP(w, r)
			return
		}

		user := user.New()
		session.Values["user"] = user.UUID
		m.Config.Log.Debug("Created new User for session: " + user.UUID)

		session.Options.MaxAge = 3600 // 1 hour
		if err := session.Save(r, w); err != nil {
			m.Config.Log.Error(err)
			write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		if err := h.SaveUser(user); err != nil {
			m.Config.Log.Error(err)
			write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) LoadSessionAndUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := helper.WithConfig(m.Config)

		session, err := m.Config.SessionStore.Get(r, "audio-gonverter")
		if err != nil {
			m.Config.Log.Error(err)
			write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		e, err := m.Config.DatabaseRepo.Exist(session.Values["user"].(string))
		if !e && err == nil {
			u := &user.User{UUID: session.Values["user"].(string), Files: []*file.File{}}
			err2 := h.SaveUser(u)
			if err2 != nil {
				m.Config.Log.Error(err)
				write(w, err2.Error(), http.StatusInternalServerError)
				return
			}

		} else if err != nil {
			m.Config.Log.Error(err)
			write(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := h.LoadUser(session.Values["user"].(string))
		if err != nil {
			m.Config.Log.Error(err)
			write(w, "Session error. Cleaning the cache may solve the issue", http.StatusInternalServerError)
			return
		}
		sr := r.WithContext(user.ToContext(r.Context()))
		next.ServeHTTP(w, sr)
	})
}

func (m *Middleware) StatusCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := user.FromRequest(r)
		handler := Handler{Config: m.Config}

		if user.IsConverting {
			handler.render(w, "status.page.gohtml", TemplateData{Messages: []string{"Converting"}, Commit: m.Config.Env["COMMIT"]})
			return
		}
		if user.IsUploading {
			handler.render(w, "status.page.gohtml", TemplateData{Messages: []string{"Uploading"}, Commit: m.Config.Env["COMMIT"]})
			return
		}

		next.ServeHTTP(w, r)
	})
}
