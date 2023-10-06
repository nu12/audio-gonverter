package main

import (
	"log"
	"net/http"

	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/ffmpeg"
	"github.com/nu12/audio-gonverter/internal/helper"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
	"github.com/nu12/audio-gonverter/internal/web"
)

func main() {
	q := &rabbitmq.RabbitQueue{}
	app := config.
		New("./cmd/templates/", "./cmd/static/").
		WithLog(logging.NewLogger()).
		LoadEnv([]string{
			"WEB_ENABLED",
			"WORKER_ENABLED",
			"REDIS_HOST",
			"REDIS_PORT",
			"QUEUE_CONNECTION_STRING",
			"COMMIT",
		}).
		LoadConfigs()

	app.WithDatabaseRepo(database.NewRedis(app.Env["REDIS_HOST"], app.Env["REDIS_PORT"], "")).
		WithConvertionToolRepo(&ffmpeg.Ffmpeg{
			InputPath:  app.OriginalPath,
			OutputPath: app.ConvertedPath,
		}).
		ConnectQueueRepo(q, app.Env["QUEUE_CONNECTION_STRING"])

	if app.Err() != nil {
		log.Fatal(app.Err())
	}
	defer q.Connection.Close()
	defer q.Channel.Close()

	c := make(chan error, 1)
	helper := &helper.Helper{Config: app}

	if app.Env["WEB_ENABLED"] == "true" {
		s := &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: web.Routes(app),
		}
		go helper.StartWeb(c, s)
	}

	if app.Env["WORKER_ENABLED"] == "true" {
		go helper.StartWorker(c)
	}

	// Panic with error
	for err := range c {
		log.Fatal(err)
	}
}
