package helper

import (
	"testing"
	"time"

	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/ffmpeg"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
)

func TestStartWorker(t *testing.T) {
	c := make(chan error)
	app := &config.Config{
		QueueRepo:          &rabbitmq.QueueMock{},
		DatabaseRepo:       &database.MockDB{},
		ConvertionToolRepo: &ffmpeg.FfmpegMock{},
		Log:                &logging.Log{},
	}

	t.Run("Start Worker", func(t *testing.T) {

		go WithConfig(app).StartWorker(c)

		select {
		case err := <-c:
			t.Errorf("Unexpected error: %s", err)
		case <-time.After(1 * time.Second):
			// No error occurred, the test passes
		}
	})

}
