package main

import (
	"html/template"
	"net/http"

	"github.com/nu12/audio-gonverter/internal/model"
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
	// TODO: send to queue

	// Proof of concept
	// err := ffmpeg.Input("/tmp/uuid.mp3").
	// 	Output("/tmp/uuid.ogg").
	// 	OverWriteOutput().ErrorToStdOut().Run()

	// if err != nil {
	// 	log.Error(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Status(w http.ResponseWriter, r *http.Request) {
	log.Debug("Status page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
