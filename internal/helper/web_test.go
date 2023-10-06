package helper

import (
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/logging"
	"github.com/nu12/audio-gonverter/internal/user"
)

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
	app := &config.Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
		Log:          &logging.Log{},
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	os.Setenv("MAX_FILES_PER_USER", "10")
	os.Setenv("MAX_FILE_SIZE", "100000")
	os.Setenv("MAX_TOTAL_SIZE_PER_USER", "1000000")
	os.Setenv("ORIGINAL_FILE_EXTENTION", "mp3")
	os.Setenv("TARGET_FILE_EXTENTION", "ogg,aac")
	os.Setenv("ORIGINAL_FILES_PATH", "/tmp")
	os.Setenv("CONVERTED_FILES_PATH", "/tmp")
	defer os.Unsetenv("MAX_FILES_PER_USER")
	defer os.Unsetenv("MAX_FILE_SIZE")
	defer os.Unsetenv("MAX_TOTAL_SIZE_PER_USER")
	defer os.Unsetenv("ORIGINAL_FILE_EXTENTION")
	defer os.Unsetenv("TARGET_FILE_EXTENTION")
	defer os.Unsetenv("ORIGINAL_FILES_PATH")
	defer os.Unsetenv("CONVERTED_FILES_PATH")

	t.Run("Start Web Service", func(t *testing.T) {
		testServer := &TestServer{}
		h := &Helper{}
		go h.WithConfig(app).StartWeb(c, testServer)

		select {
		case err := <-c:
			t.Errorf("Unexpected error: %s", err)
		case <-time.After(1 * time.Second):
			// No error occurred, the test passes
		}
	})

}

func TestAddFile(t *testing.T) {
	app := &config.Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
		OriginalPath: "/tmp",
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	user := user.New()

	file, err := file.NewFile("test.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	h := &Helper{}
	h.WithConfig(app).addFile(user, file)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestAddFilesAndSave(t *testing.T) {

	app := &config.Config{
		DatabaseRepo: &database.MockDB{},
		Env:          map[string]string{},
		Log:          &logging.Log{},
	}
	os.Setenv("SESSION_KEY", "test-key")
	defer os.Unsetenv("SESSION_KEY")

	user := user.New()
	user.IsUploading = true

	files, err := file.FilesFromForm([]*multipart.FileHeader{
		{Filename: "file.mp3"},
	})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	h := &Helper{}
	h.WithConfig(app).AddFilesAndSave(user, files)

	if user.IsUploading {
		t.Errorf("Error uploading files")
	}
}
