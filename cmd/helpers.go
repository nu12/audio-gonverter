package main

import (
	"errors"
	"fmt"
	"io"
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

func (app *Config) addFile(user *model.User, header model.RawFile) error {
	mf, err := header.Open()
	if err != nil {
		return err
	}
	defer mf.Close()
	of, err := os.OpenFile("/tmp/"+header.Filename(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer of.Close()

	bytes, err := io.Copy(of, mf)
	if err != nil {
		return err
	}

	file, err := model.NewFile(header.Filename(), bytes)
	if err != nil {
		return err
	}

	if err := user.AddFile(file); err != nil {
		return err
	}

	return nil
}

func (app *Config) addFilesAndSave(user *model.User, files []model.RawFile) {
	for _, header := range files {
		//TODO: validate max size
		//TODO: validate max files per user
		//TODO: validate file type

		if err := app.addFile(user, header); err != nil {
			log.Warning("Warning adding file")
		}
	}
	user.IsUploading = false
	if err := app.saveUser( /* Repo, */ user); err != nil {
		log.Warning("Error saving user")
	}
}

func (app *Config) saveUser(user *model.User) error {
	return app.DatabaseRepo.Save(user)
}

func (app *Config) loadUser(id string) *model.User {
	return app.DatabaseRepo.Load(id)
}
