package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/logging"
)

func TestRoutes(t *testing.T) {
	testApp := &config.Config{
		TemplatesPath:   "./../../cmd/templates/",
		StaticFilesPath: "./../../cmd/static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
		Log:             &logging.Log{},
	}
	h := Routes(testApp)

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
func TestStaticFiles(t *testing.T) {
	testApp := &config.Config{
		TemplatesPath:   "./../../cmd/templates/",
		StaticFilesPath: "./../../cmd/static/",
		SessionStore:    sessions.NewCookieStore([]byte("test-key")),
		DatabaseRepo:    &database.MockDB{},
		Log:             &logging.Log{},
	}
	h := Routes(testApp)

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
