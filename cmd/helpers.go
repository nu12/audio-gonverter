package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/model"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (app *Config) loadEnv(required []string) error {
	for _, env := range required {
		val, isSet := os.LookupEnv(env)
		if !isSet {
			return fmt.Errorf("Variable %s was not found", env)
		}
		app.Env[env] = val
	}
	return nil
}

// For testability
type CustomHTTPServer interface {
	ListenAndServe() error
}

func (app *Config) loadConfigs() error {
	tempApp := &Config{Env: map[string]string{}}
	err := tempApp.loadEnv([]string{
		"MAX_FILES_PER_USER",
		"MAX_FILE_SIZE",
		"MAX_TOTAL_SIZE_PER_USER",
		"ORIGINAL_FILE_EXTENTION",
		"TARGET_FILE_EXTENTION",
	})
	if err != nil {
		return err
	}

	app.MaxFilesPerUser, err = strconv.Atoi(tempApp.Env["MAX_FILES_PER_USER"])
	if err != nil {
		return err
	}

	app.MaxFileSize, err = strconv.Atoi(tempApp.Env["MAX_FILE_SIZE"])
	if err != nil {
		return err
	}
	app.MaxTotalSizePerUser, err = strconv.Atoi(tempApp.Env["MAX_TOTAL_SIZE_PER_USER"])
	if err != nil {
		return err
	}

	app.OriginFileExtention = strings.Split(tempApp.Env["ORIGINAL_FILE_EXTENTION"], ",")
	app.TargetFileExtention = strings.Split(tempApp.Env["TARGET_FILE_EXTENTION"], ",")

	return nil
}

func (app *Config) startWeb(c chan<- error, s CustomHTTPServer) {
	log.Info("Starting Web service")

	if err := app.loadEnv([]string{"SESSION_KEY"}); err != nil {
		c <- err
	}
	app.SessionStore = sessions.NewCookieStore([]byte(app.Env["SESSION_KEY"]))

	if err := app.loadConfigs(); err != nil {
		c <- err
	}

	if err := s.ListenAndServe(); err != nil {
		c <- err
	}
}

func (app *Config) startWorker(c chan<- error) {
	log.Info("Starting Worker service")

	for {
		msg, err := app.QueueRepo.Pull()
		if err != nil {
			c <- err
		}
		decoded, err := rabbitmq.Decode(msg)
		if err != nil {
			log.Warning("Cannot decode de message: " + msg)
			continue
		}
		user, err := app.loadUser(decoded.UserUUID)
		if err != nil {
			log.Warning("Cannot retrieve user: " + decoded.UserUUID)
			continue
		}
		if err := app.convert(user, decoded.Format, decoded.Kbps); err != nil {
			// TODO: message to the user
			log.Warning("Error converting file")
		}
		user.IsConverting = false
		if err := app.saveUser(user); err != nil {
			log.Warning("Error saving user")
			continue
		}
	}
}

func (app *Config) convert(user *model.User, format, kpbs string) error {
	for _, file := range user.Files {

		convertedId := model.GenerateUUID()
		convertedName := file.Prefix() + "." + format

		if err := os.Mkdir(CONVERTED_PATH+convertedId, 0777); err != nil {
			log.Error(err)
			return err
		}

		err := ffmpeg.Input(ORIGINAL_PATH+file.OriginalId+"/"+file.OriginalName).
			Output(CONVERTED_PATH+convertedId+"/"+convertedName, ffmpeg.KwArgs{"b:a": kpbs + "k"}).
			// OverWriteOutput().
			// ErrorToStdOut().
			Run()

		if err != nil {
			log.Warning("Could not convert file: " + err.Error())
			// TODO: needs refactoring
			err2 := user.RemoveFile(file.OriginalId)
			if err2 != nil {
				log.Warning("Error removing file: " + err2.Error())
			}
			continue
		}

		file.ConvertedName = convertedName
		file.ConvertedId = convertedId
		file.IsConverted = true
	}
	return nil
}

func (app *Config) addFile(user *model.User, file *model.File) error {

	if err := file.SaveToDisk(ORIGINAL_PATH); err != nil {
		return err
	}
	if err := user.AddFile(file); err != nil {
		return err
	}
	if err := app.saveUser(user); err != nil {
		return err
	}

	return nil
}

func (app *Config) addFilesAndSave(user *model.User, files []*model.File) {
	for _, file := range files {
		file.ValidateMaxFilesPerUser(user, app.MaxFilesPerUser)
		file.ValidateMaxSize(app.MaxFileSize)
		file.ValidateMaxSizePerUser(user, app.MaxTotalSizePerUser)
		file.ValidateFileExtention(app.OriginFileExtention)
		if message, valid := file.GetValidity(); !valid {
			log.Debug(message)
			// TODO: message to the user
			continue
		}

		if err := app.addFile(user, file); err != nil {
			log.Warning(err.Error())
		}
	}
	user.IsUploading = false
	if err := app.saveUser(user); err != nil {
		log.Warning(err.Error())
	}
}

func (app *Config) saveUser(user *model.User) error {
	return app.DatabaseRepo.Save(user)
}

func (app *Config) loadUser(id string) (*model.User, error) {
	return app.DatabaseRepo.Load(id)
}

func (app *Config) AddFlash(w http.ResponseWriter, r *http.Request, msg string) {
	s, err := app.SessionStore.Get(r, "audio-gonverter")
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.AddFlash(msg)
	err = s.Save(r, w)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) GetFlash(w http.ResponseWriter, r *http.Request) string {
	s, err := app.SessionStore.Get(r, "audio-gonverter")
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	f := s.Flashes()
	if len(f) == 0 {
		return "Welcome to audio-gonverter"
	}

	msg := f[0].(string)
	err = s.Save(r, w)
	if err != nil {
		app.write(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return msg
}
func sliceToString(s []string) string {
	ps := []string{}
	for _, f := range s {
		ps = append(ps, "."+f)
	}
	return strings.Join(ps, ",")
}
