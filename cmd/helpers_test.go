package main

import (
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/model"
	"github.com/nu12/audio-gonverter/internal/rabbitmq"
)

// Test loadEnv
func TestLoadEnv(t *testing.T) {
	testConfig := &Config{
		Env: map[string]string{},
	}

	if err := testConfig.loadEnv([]string{"EMPTY"}); err == nil {
		t.Errorf("Empty env variable should return an error, but error is nil")
	}

	env := "EXISTS"
	val := "yes"
	os.Setenv(env, val)
	if err := testConfig.loadEnv([]string{env}); err != nil {
		t.Errorf("Existing env variable should not return an error, but got %s", err)
	}

	if testConfig.Env[env] != val {
		t.Errorf("Config should contain %s with value %s", env, val)
	}
}

// Test startWeb
type TestServer struct {
	Addr    string
	Handler http.Handler
}

func (*TestServer) ListenAndServe() error {
	return nil
}

func TestStartWeb(t *testing.T) {
	c := make(chan error)
	app := Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	t.Run("Start Web Service", func(t *testing.T) {
		testServer := &TestServer{}
		go app.startWeb(c, testServer)

		select {
		case err := <-c:
			t.Errorf("Unexpected error: %s", err)
		case <-time.After(1 * time.Second):
			// No error occurred, the test passes
		}
	})

}

func TestStartWorker(t *testing.T) {
	c := make(chan error)
	app := Config{
		QueueRepo: &rabbitmq.QueueMock{},
	}

	t.Run("Start Worker", func(t *testing.T) {

		go app.startWorker(c)

		select {
		case err := <-c:
			t.Errorf("Unexpected error: %s", err)
		case <-time.After(1 * time.Second):
			// No error occurred, the test passes
		}
	})

}

func TestAddFile(t *testing.T) {
	app := Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	user := model.NewUser()

	file, err := model.NewFile("test.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	err = app.addFile(&user, file)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestAddFilesAndSave(t *testing.T) {

	app := Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	user := model.NewUser()
	user.IsUploading = true

	files, err := model.FilesFromForm([]*multipart.FileHeader{
		{Filename: "file.mp3"},
	})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	app.addFilesAndSave(&user, files)

	if user.IsUploading {
		t.Errorf("Error uploading files")
	}
}

func TestFlash(t *testing.T) {
	testApp := &Config{
		SessionStore: sessions.NewCookieStore([]byte("test")),
	}

	w := TestResponseWriter{}
	r := &http.Request{}

	expected := "Test message"

	testApp.AddFlash(w, r, expected)
	got := testApp.GetFlash(w, r)

	if expected != got {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
