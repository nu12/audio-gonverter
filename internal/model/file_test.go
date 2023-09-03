package model

import (
	"reflect"
	"regexp"
	"testing"
)

var prefix = "new_file"
var extention = "mp3"
var orignalName = prefix + "." + extention

var validUUID = regexp.MustCompile(`\S{8}-\S{4}-\S{4}-\S{4}-\S{12}`)

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
