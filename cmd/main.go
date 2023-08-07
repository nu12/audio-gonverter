package main

import (
	"log"
	"net/http"
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

func main() {
	log.Println("audio-gonverter starts here")
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
		log.Panicln(err)
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
		log.Panicln(err)
	}
}
