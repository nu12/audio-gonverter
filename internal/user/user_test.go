package user

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"

	"github.com/nu12/audio-gonverter/internal/file"
)

var validUUID = regexp.MustCompile(`\S{8}-\S{4}-\S{4}-\S{4}-\S{12}`)

func TestNewUser(t *testing.T) {

	u := New()

	if reflect.TypeOf(u) != reflect.TypeOf(&User{}) {
		t.Errorf("User type doesn't match")
	}

	if !validUUID.Match([]byte(u.UUID)) {
		t.Errorf("File UUID is not valid. Got %s", u.UUID)
	}

	if u.IsConverting != false {
		t.Errorf("User should be initialized with IsConverting equals to false")
	}

	if u.IsUploading != false {
		t.Errorf("User should be initialized with IsUploading equals to false")
	}

	if len(u.Files) > 0 {
		t.Errorf("User should be initialized with empty file list")
	}
}

func TestAddAndRemoveFile(t *testing.T) {
	u := New()
	f1, _ := file.NewFile("test1.mp3")
	f2, _ := file.NewFile("test2.mp3")

	if err := u.AddFile(f1).Err(); err != nil {
		t.Errorf(err.Error())
	}

	if len(u.Files) != 1 {
		t.Errorf("User should have 1 file")
	}

	if err := u.AddFile(f2).Err(); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 2 {
		t.Errorf("User should have 2 file")
	}

	if err := u.RemoveFile(f2.OriginalId).Err(); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 1 {
		t.Errorf("User should have 1 file after removal")
	}

	if u.Files[0].OriginalId != f1.OriginalId {
		t.Errorf("File name doesn't match")
	}

	if err := u.ClearFiles().Err(); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 0 {
		t.Errorf("User shouldn't have any file")
	}
}

func TestMessages(t *testing.T) {
	message := "Test"
	user := New()
	user.AddMessage(message)
	if len(user.Messages) != 1 {
		t.Errorf("Expected user to have 1 message, got %d", len(user.Messages))
	}

	if got := user.GetMessages()[0]; got != message {
		t.Errorf("Expected message %s, got %s", message, got)
	}

	if got := user.GetMessages()[0]; got != "Welcome to audio-gonverter!" {
		t.Errorf("Expected default message, got %s", got)
	}
}

func TestFromRequest(t *testing.T) {
	var u *User
	r := &http.Request{}
	u = FromRequest(r)

	if u.Err() == nil {
		t.Errorf("Expected error, got nil")
	}

	sr := r.WithContext(New().ToContext(r.Context()))
	u = FromRequest(sr)

	if u.Err() != nil {
		t.Errorf("Unexpected error")
	}
}
