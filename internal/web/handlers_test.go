package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
)

type TestResponseWriter struct{}

func (w TestResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w TestResponseWriter) Write(m []byte) (int, error) {
	return len(m), nil
}
func (w TestResponseWriter) WriteHeader(statusCode int) {}

func TestWriter(t *testing.T) {
	var w TestResponseWriter
	//testApp := &Config{}
	write(w, "Test message", http.StatusOK)
}

func TestUpload(t *testing.T) {
	app := &config.Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
	}
	h := Routes(app, &logging.Log{})

	// https://stackoverflow.com/questions/43904974/testing-go-http-request-formfile
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		_, err := writer.CreateFormFile("files", "someaudio.mp3")
		if err != nil {
			t.Error(err)
		}

	}()

	req := httptest.NewRequest("POST", "/upload", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, but got %d", http.StatusSeeOther, rr.Code)
	}

}

func TestConvert(t *testing.T) {
	queue := &rabbitmq.QueueMock{}

	testApp := &config.Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
		QueueRepo:       queue,
	}
	h := Routes(testApp, &logging.Log{})

	form := url.Values{}
	form.Add("format", "ogg")
	form.Add("kbps", "64")
	req := httptest.NewRequest("POST", "/convert", strings.NewReader(form.Encode()))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, but got %d", http.StatusSeeOther, rr.Code)
	}

	if queue.Count != 1 {
		t.Errorf("Expected 1 message in the queue, got %d", queue.Count)
	}

}
