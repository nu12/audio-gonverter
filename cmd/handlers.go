package main

import (
	"html/template"
	"net/http"

	"github.com/nu12/audio-gonverter/internal/model"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
)

type TemplateData struct {
	Files  []model.File
	Commit string
}

func (app *Config) Home(w http.ResponseWriter, r *http.Request) {
	log.Debug("Home page")
	user := r.Context().Value(userID("user")).(*model.User)

	td := TemplateData{
		Commit: app.Env["COMMIT"],
		Files:  user.Files,
	}

	t, err := template.New("index.page.gohtml").ParseFiles(app.TemplatesPath + "index.page.gohtml")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, td); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *Config) Upload(w http.ResponseWriter, r *http.Request) {
	log.Debug("Upload")

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := r.Context().Value(userID("user")).(*model.User)
	user.IsUploading = true
	if err := app.saveUser(user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := model.FilesFromForm(r.MultipartForm.File["files"])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go app.addFilesAndSave(user, files)

	// TODO: write message to display and redirect to index
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Convert(w http.ResponseWriter, r *http.Request) {
	log.Debug("Convert")

	err := r.ParseForm()

	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := r.Context().Value(userID("user")).(*model.User)
	message := rabbitmq.Message{
		UserUUID: user.UUID,
		Format:   r.PostForm.Get("format"),
		Kbps:     r.PostForm.Get("kbps"),
	}

	encoded, err := rabbitmq.Encode(message)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.QueueRepo.Push(encoded)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Sent message: " + encoded)

	user.IsConverting = true
	if err := app.saveUser(user); err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Status(w http.ResponseWriter, r *http.Request) {
	log.Debug("Status page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
