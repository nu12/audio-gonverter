package main

import (
	"io"
	"net/http"
	"os"

	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/repository"
	"github.com/nu12/audio-gonverter/internal/user"
)

type TemplateData struct {
	Files      []*file.File
	FilesCount int
	Commit     string
	Messages   []string
	Accepted   string
	Formats    []string
}

func (app *Config) Home(w http.ResponseWriter, r *http.Request) {
	log.Debug("Home page")
	user := user.FromRequest(r)

	td := TemplateData{
		Commit:     app.Env["COMMIT"],
		Files:      user.Files,
		FilesCount: len(user.Files),
		Messages:   app.GetFlash(user),
		Accepted:   sliceToString(app.OriginFileExtention),
		Formats:    app.TargetFileExtention,
	}

	app.render(w, "index.page.gohtml", td)
}

func (app *Config) Upload(w http.ResponseWriter, r *http.Request) {
	log.Debug("Upload")

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := user.FromRequest(r)
	user.IsUploading = true
	if err := app.saveUser(user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := file.FilesFromForm(r.MultipartForm.File["files"])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go app.addFilesAndSave(user, files)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Convert(w http.ResponseWriter, r *http.Request) {
	log.Debug("Convert")

	err := r.ParseForm()

	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := user.FromRequest(r)
	message := repository.QueueMessage{
		UserUUID: user.UUID,
		Format:   r.PostForm.Get("format"),
		Kbps:     r.PostForm.Get("kbps"),
	}

	encoded, err := app.QueueRepo.Encode(message)
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

func (app *Config) Delete(w http.ResponseWriter, r *http.Request) {
	log.Debug("Delete page")
	user := user.FromRequest(r)
	uuid := r.URL.Query().Get("uuid")
	if err := user.RemoveFile(uuid); err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := app.saveUser(user); err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.AddFlash(user, "Deleted file")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) DeleteAll(w http.ResponseWriter, r *http.Request) {
	log.Debug("DeleteAll page")
	user := user.FromRequest(r)
	if err := user.ClearFiles(); err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := app.saveUser(user); err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.AddFlash(user, "Deleted all files")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Download(w http.ResponseWriter, r *http.Request) {
	log.Debug("Download page")

	uuid := r.URL.Query().Get("uuid")
	dir, err := os.Open(app.ConvertedPath + "/" + uuid)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := dir.Readdirnames(-1)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.Open(app.ConvertedPath + "/" + uuid + "/" + files[0])
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b = make([]byte, 1024)
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment;filename="+files[0])

	for {
		_, err = f.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			app.write(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			app.write(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
