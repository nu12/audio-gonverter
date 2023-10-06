package web

import (
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/helper"
	"github.com/nu12/audio-gonverter/internal/repository"
	"github.com/nu12/audio-gonverter/internal/user"
)

type Handler struct {
	Config *config.Config
}

type TemplateData struct {
	Files      []*file.File
	FilesCount int
	Commit     string
	Messages   []string
	Accepted   string
	Formats    []string
}

func write(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) render(w http.ResponseWriter, templateName string, td TemplateData) {
	t, err := template.New(templateName).ParseFiles(handler.Config.TemplatesPath+"base.layout.gohtml", handler.Config.TemplatesPath+templateName)
	if err != nil {
		handler.Config.Log.Error(err)
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, td); err != nil {
		handler.Config.Log.Error(err)
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) Home(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("Home page")
	user := user.FromRequest(r)
	h := &helper.Helper{}
	h.WithConfig(handler.Config)

	td := TemplateData{
		Commit:     handler.Config.Env["COMMIT"],
		Files:      user.Files,
		FilesCount: len(user.Files),
		Messages:   h.GetFlash(user),
		Accepted:   helper.SliceToString(handler.Config.OriginFileExtention),
		Formats:    handler.Config.TargetFileExtention,
	}

	handler.render(w, "index.page.gohtml", td)
}

func (handler *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("Upload")

	h := &helper.Helper{}
	h.WithConfig(handler.Config)

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		handler.Config.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := user.FromRequest(r)
	user.IsUploading = true
	if err := h.SaveUser(user); err != nil {
		handler.Config.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := file.FilesFromForm(r.MultipartForm.File["files"])
	if err != nil {
		handler.Config.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.AddFilesAndSave(user, files)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("Convert")
	h := &helper.Helper{}
	h.WithConfig(handler.Config)

	err := r.ParseForm()

	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := user.FromRequest(r)
	message := repository.QueueMessage{
		UserUUID: user.UUID,
		Format:   r.PostForm.Get("format"),
		Kbps:     r.PostForm.Get("kbps"),
	}

	encoded, err := handler.Config.QueueRepo.Encode(message)
	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.Config.QueueRepo.Push(encoded)
	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	handler.Config.Log.Debug("Sent message: " + encoded)

	user.IsConverting = true
	if err := h.SaveUser(user); err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("Delete page")
	h := &helper.Helper{}
	h.WithConfig(handler.Config)

	uuid := r.URL.Query().Get("uuid")
	user := user.FromRequest(r)
	if err := user.RemoveFile(uuid).Err(); err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.SaveUser(user); err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.AddFlash(user, "Deleted file")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("DeleteAll page")
	h := &helper.Helper{}
	h.WithConfig(handler.Config)
	user := user.FromRequest(r).ClearFiles()

	if err := h.SaveUser(user); err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.AddFlash(user, "Deleted all files")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) Download(w http.ResponseWriter, r *http.Request) {
	handler.Config.Log.Debug("Download page")

	uuid := r.URL.Query().Get("uuid")
	dir, err := os.Open(handler.Config.ConvertedPath + "/" + uuid)
	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := dir.Readdirnames(-1)
	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.Open(handler.Config.ConvertedPath + "/" + uuid + "/" + files[0])
	if err != nil {
		write(w, err.Error(), http.StatusInternalServerError)
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
			write(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			write(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
