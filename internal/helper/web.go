package helper

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/user"
)

// For testability
type CustomHTTPServer interface {
	ListenAndServe() error
}

func (h *Helper) StartWeb(c chan<- error, s CustomHTTPServer) {
	h.Config.Log.Info("Starting Web service")
	h.Config.LoadEnv([]string{"SESSION_KEY"})
	h.Config.SessionStore = sessions.NewCookieStore([]byte(h.Config.Env["SESSION_KEY"]))

	if h.Config.Err() != nil {
		c <- h.Config.Err()
	}

	if err := s.ListenAndServe(); err != nil {
		c <- err
	}
}

func (h *Helper) AddFilesAndSave(user *user.User, files []*file.File) {
	for _, file := range files {
		file.ValidateMaxFilesPerUser(user.Files, h.Config.MaxFilesPerUser)
		file.ValidateMaxSize(h.Config.MaxFileSize)
		file.ValidateMaxSizePerUser(user.Files, h.Config.MaxTotalSizePerUser)
		file.ValidateFileExtention(h.Config.OriginFileExtention)
		if message, valid := file.GetValidity(); !valid {
			h.Config.Log.Debug(message)
			user.AddMessage(fmt.Sprintf("File %s: %s.", file.OriginalName, message))
			continue
		}

		if err := h.addFile(user, file); err != nil {
			h.Config.Log.Warning(err.Error())
		}
	}
	user.IsUploading = false
	if err := h.SaveUser(user); err != nil {
		h.Config.Log.Warning(err.Error())
	}
}

func (h *Helper) addFile(user *user.User, file *file.File) error {

	if err := file.SaveToDisk(h.Config.OriginalPath); err != nil {
		return err
	}
	if err := user.AddFile(file).Err(); err != nil {
		return err
	}
	if err := h.SaveUser(user); err != nil {
		return err
	}

	return nil
}
