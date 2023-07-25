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
	Env map[string]string
}

// For testability
type CustomHTTPServer interface {
	ListenAndServe() error
}

func main() {
	log.Println("audio-gonverter starts here")
	app := Config{
		Env: map[string]string{},
	}

	err := app.loadEnv([]string{"WEB_ENABLED", "WORKER_ENABLED"})
	if err != nil {
		log.Panicln(err)
	}

	c := make(chan error, 1)

	if app.Env["WEB_ENABLED"] == "true" {
		go app.startWeb(c, &http.Server{
			Addr: "0.0.0.0:8080",
		})
	}

	if app.Env["WORKER_ENABLED"] == "true" {
		go app.startWorker(c)
	}

	// Panic with error
	for err := range c {
		log.Panicln(err)
	}
}
