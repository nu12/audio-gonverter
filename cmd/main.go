package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/ffmpeg"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
	"github.com/nu12/audio-gonverter/internal/repository"
)

const (
	ORIGINAL_PATH  = "/tmp/original/"
	CONVERTED_PATH = "/tmp/converted/"
)

type Config struct {
	TemplatesPath       string
	StaticFilesPath     string
	SessionStore        *sessions.CookieStore
	DatabaseRepo        repository.DatabaseRepository
	QueueRepo           repository.QueueRepository
	ConvertionToolRepo  repository.ConvertionToolRepo
	Env                 map[string]string
	MaxFilesPerUser     int
	MaxFileSize         int
	MaxTotalSizePerUser int
	OriginFileExtention []string
	TargetFileExtention []string
}

var log = logging.NewLogger()

func main() {
	log.Info("Starting audio-gonverter")

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
		"QUEUE_CONNECTION_STRING",
		"COMMIT",
	})
	if err != nil {
		log.Fatal(err)
	}

	app.DatabaseRepo = database.NewRedis(app.Env["REDIS_HOST"], app.Env["REDIS_PORT"], "")
	app.ConvertionToolRepo = &ffmpeg.Ffmpeg{
		InputPath:  ORIGINAL_PATH,
		OutputPath: CONVERTED_PATH,
	}

	q := &rabbitmq.RabbitQueue{}
	err = errors.New("No queue available")
	for i := 1; i <= 10; i++ {
		q, err = rabbitmq.Connect(app.Env["QUEUE_CONNECTION_STRING"])
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			log.Debug("Queue not ready, waiting...")
			continue
		}
		break
	}

	if err != nil {
		log.Fatal(err)
	}
	app.QueueRepo = q
	defer q.Connection.Close()
	defer q.Channel.Close()

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
