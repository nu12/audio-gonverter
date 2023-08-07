package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	testApp := &Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
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

func TestStaticFiles(t *testing.T) {
	testApp := &Config{
		TemplatesPath:   "./templates/",
		StaticFilesPath: "./static/",
	}
	h := testApp.routes()

	for _, file := range []string{
		"/static/javascript.js",
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
