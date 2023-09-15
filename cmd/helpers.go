package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/model"
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

	c <- errors.New("Worker is not implemented")
}

func (app *Config) addFile(user *model.User, file *model.File) error {

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
