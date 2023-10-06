package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
	"github.com/nu12/audio-gonverter/internal/repository"
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
	OriginalPath        string
	ConvertedPath       string
	Log                 *logging.Log

	err error
}

func New(templatesPath, staticFilesPath string) *Config {
	return &Config{
		TemplatesPath:   templatesPath,
		StaticFilesPath: staticFilesPath,
		Env:             map[string]string{},
		err:             nil,
	}
}

func (app *Config) Err() error {
	return app.err
}

func (app *Config) LoadEnv(required []string) *Config {
	for _, env := range required {
		if app.err != nil {
			return app
		}
		val, isSet := os.LookupEnv(env)
		if !isSet {
			app.err = fmt.Errorf("Variable %s was not found", env)
		}
		app.Env[env] = val
	}
	return app
}

func (app *Config) LoadConfigs() *Config {
	if app.err != nil {
		return app
	}

	tempApp := &Config{Env: map[string]string{}}
	tempApp.LoadEnv([]string{
		"MAX_FILES_PER_USER",
		"MAX_FILE_SIZE",
		"MAX_TOTAL_SIZE_PER_USER",
		"ORIGINAL_FILE_EXTENTION",
		"TARGET_FILE_EXTENTION",
		"ORIGINAL_FILES_PATH",
		"CONVERTED_FILES_PATH",
	})

	app.MaxFilesPerUser, app.err = strconv.Atoi(tempApp.Env["MAX_FILES_PER_USER"])
	if app.err != nil {
		return app
	}

	app.MaxFileSize, app.err = strconv.Atoi(tempApp.Env["MAX_FILE_SIZE"])
	if app.err != nil {
		return app
	}
	app.MaxTotalSizePerUser, app.err = strconv.Atoi(tempApp.Env["MAX_TOTAL_SIZE_PER_USER"])
	if app.err != nil {
		return app
	}

	app.OriginalPath = tempApp.Env["ORIGINAL_FILES_PATH"]
	app.ConvertedPath = tempApp.Env["CONVERTED_FILES_PATH"]

	app.OriginFileExtention = strings.Split(tempApp.Env["ORIGINAL_FILE_EXTENTION"], ",")
	app.TargetFileExtention = strings.Split(tempApp.Env["TARGET_FILE_EXTENTION"], ",")

	return app
}

func (app *Config) WithDatabaseRepo(r repository.DatabaseRepository) *Config {
	if app.err != nil {
		return app
	}

	app.DatabaseRepo = r
	return app
}

func (app *Config) WithConvertionToolRepo(r repository.ConvertionToolRepo) *Config {
	if app.err != nil {
		return app
	}

	app.ConvertionToolRepo = r
	return app
}

func (app *Config) WithLog(l *logging.Log) *Config {
	if app.err != nil {
		return app
	}

	app.Log = l
	return app
}

func (app *Config) ConnectQueueRepo(q repository.QueueRepository, conn string) *Config {
	if app.err != nil {
		return app
	}

	err := errors.New("No queue available")

	for i := 1; i <= 10; i++ {
		q, err = rabbitmq.Connect(conn)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		break
	}

	app.QueueRepo = q
	app.err = err

	return app
}
