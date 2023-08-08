package main

import (
	"net/http"

	"github.com/nu12/audio-gonverter/internal/logging"
)

const (
	ORIGINAL_PATH  = "/tmp/original/"
	CONVERTED_PATH = "/tmp/converted/"
)

type Config struct {
	TemplatesPath   string
	StaticFilesPath string
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
		"COMMIT",
	})
	if err != nil {
		log.Fatal(err)
	}

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
