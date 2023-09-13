package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/database"
)

func TestRoutes(t *testing.T) {
	testApp := &Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
	}
	h := testApp.routes()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

}

func TestUpload(t *testing.T) {
	testApp := &Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
	}
	h := testApp.routes()

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

func TestStaticFiles(t *testing.T) {
	testApp := &Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
	}
	h := testApp.routes()

	for _, file := range []string{
		"/static/css.css",
		"/static/logo.png",
	} {
		req, err := http.NewRequest("GET", file, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d for file %s, but got %d", http.StatusOK, file, rr.Code)
		}
	}
}

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
	testApp := &Config{}
	testApp.write(w, "Test message", http.StatusOK)
}
