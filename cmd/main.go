package main

import (
	"github.com/nu12/audio-gonverter/internal/logging"
)

var log logging.Log

func main() {
	log = *logging.NewLogger()

	c := make(chan error, 1)

	// Panic with error
	for err := range c {
		log.Fatal(err)
	}
}
