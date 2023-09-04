package model

import (
	"reflect"
	"testing"
)

func TestNewUseer(t *testing.T) {

	u := NewUser()

	if reflect.TypeOf(u) != reflect.TypeOf(User{}) {
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
	u := NewUser()
	f1, _ := NewFile("test1.mp3")
	f2, _ := NewFile("test2.mp3")

	if err := u.AddFile(f1); err != nil {
		t.Errorf(err.Error())
	}

	if len(u.Files) != 1 {
		t.Errorf("User should have 1 file")
	}

	if err := u.AddFile(f2); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 2 {
		t.Errorf("User should have 2 file")
	}

	if err := u.RemoveFile(f2.OriginalId); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 1 {
		t.Errorf("User should have 1 file after removal")
	}

	if u.Files[0].OriginalId != f1.OriginalId {
		t.Errorf("File name doesn't match")
	}

	if err := u.ClearFiles(); err != nil {
		t.Errorf(err.Error())
	}
	if len(u.Files) != 0 {
		t.Errorf("User shouldn't have any file")
	}
}
