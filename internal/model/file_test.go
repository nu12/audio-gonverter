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
	if getExtention(prefix) != "" {
		t.Errorf("Prefix doesn't match. Expected empty, got %s", getExtention(prefix))
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
	if f.raw.File != h {
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

func TestValidateFileExtention(t *testing.T) {
	validExtentions := []string{"mp3"}
	invalidExtention, err := NewFile("invalid.pdf")
	invalidityMessage := "Invalid file extention"
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if _, valid := invalidExtention.GetValidity(); !valid {
		t.Errorf("Expected a newly created file to be valid before validation.")
	}

	invalidExtention.ValidateFileExtention(validExtentions)
	if _, valid := invalidExtention.GetValidity(); !valid {
		t.Errorf("Expected file without raw file to be valid.")
	}
	invalidExtention.addRawFile(&multipart.FileHeader{})
	invalidExtention.ValidateFileExtention(validExtentions)
	if _, valid := invalidExtention.GetValidity(); valid {
		t.Errorf("Expected file with raw file to be invalid after validation.")
	}
	if message, _ := invalidExtention.GetValidity(); message != invalidityMessage {
		t.Errorf("Expected invalidity message to be %s, got %s.", invalidityMessage, message)
	}

	validExtention, err := NewFile("valid.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	validExtention.addRawFile(&multipart.FileHeader{})
	validExtention.ValidateFileExtention(validExtentions)
	if _, valid := validExtention.GetValidity(); !valid {
		t.Errorf("Expected file to be valid after validation.")
	}
}

func TestValidateMaxSize(t *testing.T) {
	maxSize := 10000
	invalidityMessage := "File is too big"

	big, err := NewFile("bid_file.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	big.addRawFile(&multipart.FileHeader{Size: 10001})

	if _, valid := big.GetValidity(); !valid {
		t.Errorf("Expected file with raw file to be valid before validation.")
	}
	big.ValidateMaxSize(maxSize)
	if _, valid := big.GetValidity(); valid {
		t.Errorf("Expected file with raw file to be invalid after validation.")
	}
	if message, _ := big.GetValidity(); message != invalidityMessage {
		t.Errorf("Expected invalidity message to be %s, got %s.", invalidityMessage, message)
	}

	small, err := NewFile("small_file.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	small.addRawFile(&multipart.FileHeader{Size: 9999})

	small.ValidateMaxSize(maxSize)
	if _, valid := small.GetValidity(); !valid {
		t.Errorf("Expected file with raw file to be valid after validation.")
	}
}

func TestValidateMaxSizePerUser(t *testing.T) {
	maxSizePerUser := 5001
	invalidityMessage := "User's file size limit reached"

	user := NewUser()
	file1, err := NewFile("File1.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	file1.addRawFile(&multipart.FileHeader{Size: 5000})
	file1.ValidateMaxSizePerUser(&user, maxSizePerUser)
	if _, valid := file1.GetValidity(); !valid {
		t.Errorf("Expected first file to be valid")
	}

	err = user.AddFile(file1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	file2, err := NewFile("File2.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	file2.addRawFile(&multipart.FileHeader{Size: 5000})
	file2.ValidateMaxSizePerUser(&user, maxSizePerUser)
	if _, valid := file2.GetValidity(); valid {
		t.Errorf("Expected second file to be invalid")
	}
	if message, _ := file2.GetValidity(); message != invalidityMessage {
		t.Errorf("Expected invalidity message to be %s, got %s.", invalidityMessage, message)
	}

}

func TestValidateMaxFilesPerUser(t *testing.T) {
	maxfilesPerUser := 1
	invalidityMessage := "User's max files limit reached"

	user := NewUser()
	file1, err := NewFile("File1.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	file1.addRawFile(&multipart.FileHeader{})
	file1.ValidateMaxFilesPerUser(&user, maxfilesPerUser)
	if _, valid := file1.GetValidity(); !valid {
		t.Errorf("Expected first file to be valid")
	}

	err = user.AddFile(file1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	file2, err := NewFile("File2.mp3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	file2.addRawFile(&multipart.FileHeader{})
	file2.ValidateMaxFilesPerUser(&user, maxfilesPerUser)
	if _, valid := file2.GetValidity(); valid {
		t.Errorf("Expected second file to be invalid")
	}
	if message, _ := file2.GetValidity(); message != invalidityMessage {
		t.Errorf("Expected invalidity message to be %s, got %s.", invalidityMessage, message)
	}

}

func TestMultipleValidation(t *testing.T) {

	file, err := NewFile("invalid.pdf")
	if err != nil {
		t.Errorf("Unexpected errer: %s", err)
	}
	file.addRawFile(&multipart.FileHeader{Size: 10})

	// This validation should fail
	file.ValidateFileExtention([]string{"mp3"})
	if _, valid := file.GetValidity(); valid {
		t.Errorf("Expected file to be invalid after first validation")
	}

	// This validation should pass
	file.ValidateMaxSize(100)
	if _, valid := file.GetValidity(); valid {
		t.Errorf("Expected file to be invalid after second validation")
	}
}
