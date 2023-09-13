package model

import (
	"mime/multipart"
	"os"
	"reflect"
	"testing"
)

var prefix = "new_file"
var extention = "mp3"
var orignalName = prefix + "." + extention

func TestNewFile(t *testing.T) {

	f, _ := NewFile(orignalName)
	if reflect.TypeOf(*f) != reflect.TypeOf(File{}) {
		t.Errorf("File type doesn't match")
	}

	if f.OriginalName != orignalName {
		t.Errorf("File name doesn't match. Expected %s, got %s", orignalName, f.OriginalName)
	}

	if !validUUID.Match([]byte(getPrefix(f.OriginalId))) {
		t.Errorf("File UUID is not valid. Got %s", getPrefix(f.OriginalId))
	}
}

func TestGetPrefix(t *testing.T) {
	if getPrefix(orignalName) != prefix {
		t.Errorf("Prefix doesn't match. Expected %s, got %s", prefix, getPrefix(orignalName))
	}
}

func TestGetExtention(t *testing.T) {
	if getExtention(orignalName) != extention {
		t.Errorf("Prefix doesn't match. Expected %s, got %s", extention, getExtention(orignalName))
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid := generateUUID()
	match := validUUID.Match([]byte(uuid))
	if !match {
		t.Errorf("Error creating UUID. Got %s", uuid)
	}
}

func TestAddandGetRawFile(t *testing.T) {
	h := &multipart.FileHeader{}
	f := &File{}
	f.addRawFile(h)
	if f.raw != h {
		t.Errorf("Error adding raw file")
	}

	if f.getRawFile() != h {
		t.Errorf("Error getting raw file")
	}
}

func TestEmptyRawFile(t *testing.T) {
	f, err := NewFile("test.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if f.getRawFile() != nil {
		t.Errorf("Raw file should be nil")
	}
}

func TestFilesFromForm(t *testing.T) {
	files, err := FilesFromForm([]*multipart.FileHeader{
		{
			Filename: "file.mp3",
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected files to have one element, got %d", len(files))
	}

}

func TestSaveToDisk(t *testing.T) {

	f, err := NewFile("test-file.mp3")
	if err != nil {
		t.Errorf("Error creating file %s", err)
	}

	err = f.SaveToDisk("/tmp")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if _, err := os.Stat("/tmp/" + f.OriginalId); err != nil {
		t.Errorf("Expected file to exist: %s", err)
	}
}
