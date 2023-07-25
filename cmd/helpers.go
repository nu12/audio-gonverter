package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (app *Config) loadEnv(required []string) error {
	for _, env := range required {
		val, isSet := os.LookupEnv(env)
		if !isSet {
			return errors.New(fmt.Sprintf("Variable %s was not found", env))
		}
		app.Env[env] = val
	}
	return nil
}

func (app *Config) startWeb(c chan<- error, s CustomHTTPServer) {
	log.Println("Stating Web service")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, gonverter")
	})

	err := s.ListenAndServe()
	if err != nil {
		c <- err
	}
}

func (app *Config) startWorker(c chan<- error) {
	log.Println("Stating Worker service")

	c <- errors.New("Worker is not implemented")
}
