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

	os.Setenv("MAX_FILES_PER_USER", "10")
	os.Setenv("MAX_FILE_SIZE", "100000")
	os.Setenv("MAX_TOTAL_SIZE_PER_USER", "1000000")
	os.Setenv("ORIGINAL_FILE_EXTENTION", "mp3")
	os.Setenv("TARGET_FILE_EXTENTION", "ogg,aac")
	defer os.Unsetenv("MAX_FILES_PER_USER")
	defer os.Unsetenv("MAX_FILE_SIZE")
	defer os.Unsetenv("MAX_TOTAL_SIZE_PER_USER")
	defer os.Unsetenv("ORIGINAL_FILE_EXTENTION")
	defer os.Unsetenv("TARGET_FILE_EXTENTION")

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

func TestLoadConfig(t *testing.T) {
	testApp := &Config{}
	os.Setenv("MAX_FILES_PER_USER", "10")
	os.Setenv("MAX_FILE_SIZE", "100000")
	os.Setenv("MAX_TOTAL_SIZE_PER_USER", "1000000")
	os.Setenv("ORIGINAL_FILE_EXTENTION", "mp3")
	os.Setenv("TARGET_FILE_EXTENTION", "ogg,aac")
	defer os.Unsetenv("MAX_FILES_PER_USER")
	defer os.Unsetenv("MAX_FILE_SIZE")
	defer os.Unsetenv("MAX_TOTAL_SIZE_PER_USER")
	defer os.Unsetenv("ORIGINAL_FILE_EXTENTION")
	defer os.Unsetenv("TARGET_FILE_EXTENTION")

	err := testApp.loadConfigs()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if testApp.MaxFilesPerUser != 10 {
		t.Errorf("Expected MaxFilesPerUser to be %d, got %d", 10, testApp.MaxFilesPerUser)
	}
	if testApp.MaxFileSize != 100000 {
		t.Errorf("Expected MaxSize to be %d, got %d", 100000, testApp.MaxFileSize)
	}
	if testApp.MaxTotalSizePerUser != 1000000 {
		t.Errorf("Expected MaxSizePerUser to be %d, got %d", 1000000, testApp.MaxTotalSizePerUser)
	}
	if len(testApp.OriginFileExtention) != 1 {
		t.Errorf("Expected OriginFileExtention to have %d items, it has %d", 1, len(testApp.OriginFileExtention))
	}
	if testApp.OriginFileExtention[0] != "mp3" {
		t.Errorf("Expected first OriginFileExtention to be %s, got %s", "mp3", testApp.OriginFileExtention[0])
	}
	if len(testApp.TargetFileExtention) != 2 {
		t.Errorf("Expected TargetFileExtention to bhave %d items, it has %d", 2, len(testApp.TargetFileExtention))
	}

	if testApp.TargetFileExtention[0] != "ogg" {
		t.Errorf("Expected first TargetFileExtention to be %s, got %s", "ogg", testApp.TargetFileExtention[0])
	}
	if testApp.TargetFileExtention[1] != "aac" {
		t.Errorf("Expected first TargetFileExtention to be %s, got %s", "aac", testApp.TargetFileExtention[1])
	}

}

func TestSliceToString(t *testing.T) {
	s := []string{"t1", "t2"}
	ps := sliceToString(s)
	if ps != ".t1,.t2" {
		t.Errorf("Expected %s, got %s", ".t1,.t2", ps)
	}
}
