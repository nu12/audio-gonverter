package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/user"
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
		"ORIGINAL_FILES_PATH",
		"CONVERTED_FILES_PATH",
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

	app.OriginalPath = tempApp.Env["ORIGINAL_FILES_PATH"]
	app.ConvertedPath = tempApp.Env["CONVERTED_FILES_PATH"]

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
		decoded, err := app.QueueRepo.Decode(msg)
		if err != nil {
			log.Warning("Cannot decode the message: " + msg)
			continue
		}
		user, err := app.loadUser(decoded.UserUUID)
		if err != nil {
			log.Warning("Cannot retrieve user: " + decoded.UserUUID)
			continue
		}
		if err := app.convert(user, decoded.Format, decoded.Kbps); err != nil {
			log.Warning("Error converting file")
		}
		user.IsConverting = false
		if err := app.saveUser(user); err != nil {
			log.Warning("Error saving user")
			continue
		}
	}
}

func (app *Config) convert(user *user.User, format, kpbs string) error {
	for _, file := range user.Files {

		err := app.ConvertionToolRepo.Convert(file, format, kpbs)
		if err != nil {
			log.Warning(err.Error())
			user.AddMessage(fmt.Sprintf("Error converting file %s (%s). Try again with different parameters.", file.OriginalName, err.Error()))
		}
	}
	return nil
}

func (app *Config) addFile(user *user.User, file *file.File) error {

	if err := file.SaveToDisk(app.OriginalPath); err != nil {
		return err
	}
	if err := user.AddFile(file).Err(); err != nil {
		return err
	}
	if err := app.saveUser(user); err != nil {
		return err
	}

	return nil
}

func (app *Config) addFilesAndSave(user *user.User, files []*file.File) {
	for _, file := range files {
		file.ValidateMaxFilesPerUser(user.Files, app.MaxFilesPerUser)
		file.ValidateMaxSize(app.MaxFileSize)
		file.ValidateMaxSizePerUser(user.Files, app.MaxTotalSizePerUser)
		file.ValidateFileExtention(app.OriginFileExtention)
		if message, valid := file.GetValidity(); !valid {
			log.Debug(message)
			user.AddMessage(fmt.Sprintf("File %s: %s.", file.OriginalName, message))
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

func (app *Config) saveUser(user *user.User) error {
	return app.DatabaseRepo.Save(user)
}

func (app *Config) loadUser(id string) (*user.User, error) {
	return app.DatabaseRepo.Load(id)
}

func (app *Config) AddFlash(u *user.User, msg string) {
	u.AddMessage(msg)
	if err := app.saveUser(u); err != nil {
		log.Warning("Error saving user with flash message: " + err.Error())
	}
}

func (app *Config) GetFlash(u *user.User) []string {
	if len(u.Messages) == 0 {
		return []string{"Welcome to audio-gonverter!"}
	}
	messages := u.GetMessages()
	if err := app.saveUser(u); err != nil {
		log.Warning("Error saving user without flash messages: " + err.Error())
	}
	return messages
}
func sliceToString(s []string) string {
	ps := []string{}
	for _, f := range s {
		ps = append(ps, "."+f)
	}
	return strings.Join(ps, ",")
}
