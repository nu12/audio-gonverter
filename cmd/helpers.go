package main

import (
	"fmt"
	"os"

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
			log.Warning("Error converting file")
			continue
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

		convertedId := model.GenerateUUID() + "." + format

		// TODO: configure
		err := ffmpeg.Input("/tmp/"+file.OriginalId).
			Output("/tmp/"+convertedId, ffmpeg.KwArgs{"b:a": kpbs + "k"}).
			OverWriteOutput().
			ErrorToStdOut().
			Run()

		if err != nil {
			return err
		}

		file.ConvertedName = file.Prefix() + "." + format
		file.ConvertedId = convertedId
		file.IsConverted = true
	}
	return nil
}

func (app *Config) addFile(user *model.User, file *model.File) error {

	// TODO: configure
	if err := file.SaveToDisk("/tmp"); err != nil {
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
		file.ValidateMaxFilesPerUser(user, 10)       //TODO: configuration
		file.ValidateMaxSize(10000000)               //TODO: configuration
		file.ValidateMaxSizePerUser(user, 100000000) //TODO: configuration
		file.ValidateFileExtention([]string{"mp3"})  //TODO: configuration
		if message, valid := file.GetValidity(); !valid {
			log.Debug(message)
			//TODO: add message to user
			break
		}

		if err := app.addFile(user, file); err != nil {
			log.Warning(err.Error())
		}
	}
	user.IsUploading = false
	if err := app.saveUser( /* Repo, */ user); err != nil {
		log.Warning(err.Error())
	}
}

func (app *Config) saveUser(user *model.User) error {
	return app.DatabaseRepo.Save(user)
}

func (app *Config) loadUser(id string) (*model.User, error) {
	return app.DatabaseRepo.Load(id)
}
