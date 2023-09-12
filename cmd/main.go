package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/repository"
)

const (
	ORIGINAL_PATH  = "/tmp/original/"
	CONVERTED_PATH = "/tmp/converted/"
)

type Config struct {
	TemplatesPath   string
	StaticFilesPath string
	SessionStore    *sessions.CookieStore
	DatabaseRepo    repository.DatabaseRepository
	Env             map[string]string
}

var log = logging.NewLogger()

func main() {
	log.Debug("audio-gonverter starts here")

	app := Config{
		TemplatesPath:   "./cmd/templates/",
		StaticFilesPath: "./cmd/static/",
		Env:             map[string]string{},
	}

	err := app.loadEnv([]string{
		"WEB_ENABLED",
		"WORKER_ENABLED",
		"REDIS_HOST",
		"REDIS_PORT",
		"COMMIT",
	})
	if err != nil {
		log.Fatal(err)
	}

	app.DatabaseRepo = database.NewRedis(app.Env["REDIS_HOST"], app.Env["REDIS_PORT"], "")
	c := make(chan error, 1)

	if app.Env["WEB_ENABLED"] == "true" {
		s := &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: app.routes(),
		}
		go app.startWeb(c, s)
	}

	if app.Env["WORKER_ENABLED"] == "true" {
		go app.startWorker(c)
	}

	// Panic with error
	for err := range c {
		log.Fatal(err)
	}
}
