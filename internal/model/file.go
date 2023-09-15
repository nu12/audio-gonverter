package model

import (
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Raw struct {
	File              *multipart.FileHeader
	IsValid           bool
	InvalidityMessage string
}

type File struct {
	OriginalName  string `json:"original_name"`
	ConvertedName string `json:"converted_name"`
	OriginalSize  int64  `json:"original_size"`
	ConvertedSize int64  `json:"converted_size"`
	OriginalId    string `json:"original_id"`
	ConvertedId   string `json:"converted_id"`
	IsConverted   bool   `json:"is_converted"`
	raw           Raw    `json:"-"`
}

func NewFile(OriginalName string) (*File, error) {
	return &File{
		OriginalName: OriginalName,
		OriginalId:   generateUUID() + "." + getExtention(OriginalName),
	}, nil
}

func FilesFromForm(rawFiles []*multipart.FileHeader) ([]*File, error) {
	var files = make([]*File, 0)

	for _, rawFile := range rawFiles {
		f, err := NewFile(rawFile.Filename)
		if err != nil {
			return files, err
		}

		f.addRawFile(rawFile)
		files = append(files, f)
	}
	return files, nil
}

func (f *File) SaveToDisk(path string) error {
	of, err := os.OpenFile(path+"/"+f.OriginalId, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer of.Close()

	raw := f.getRawFile()
	if raw == nil {
		return nil
	}

	mf, err := raw.Open()
	if err != nil {
		return err
	}
	defer mf.Close()

	bytes, err := io.Copy(of, mf)
	if err != nil {
		return err
	}

	f.OriginalSize = bytes
	return nil
}

func (f *File) GetValidity() (string, bool) {
	if f.getRawFile() == nil {
		return "", true
	}
	return f.raw.InvalidityMessage, f.raw.IsValid
}

func (f *File) ValidateFileExtention(validExtentions []string) {
	if !f.raw.IsValid {
		return
	}
	e := getExtention(f.OriginalName)
	isValid := false
	message := "Invalid file extention"
	for _, valid := range validExtentions {
		if e == valid {
			isValid = true
			message = ""
		}
	}
	f.raw.IsValid = isValid
	f.raw.InvalidityMessage = message
}

func (f *File) ValidateMaxSize(maxSize int) {
	if !f.raw.IsValid {
		return
	}
	isValid := false
	message := "File is too big"

	if f.raw.File.Size <= int64(maxSize) {
		isValid = true
		message = ""
	}

	f.raw.IsValid = isValid
	f.raw.InvalidityMessage = message
}

func (f *File) ValidateMaxSizePerUser(user *User, maxSizePerUser int) {
	if !f.raw.IsValid {
		return
	}
	isValid := false
	message := "User's file size limit reached"

	currentUserFileSize := 0
	for _, file := range user.Files {
		if file.getRawFile() != nil {
			currentUserFileSize += int(file.raw.File.Size)
		} else {
			currentUserFileSize += int(file.OriginalSize)
		}
	}

	if (currentUserFileSize + int(f.raw.File.Size)) <= maxSizePerUser {
		isValid = true
		message = ""
	}

	f.raw.IsValid = isValid
	f.raw.InvalidityMessage = message
}

func (f *File) ValidateMaxFilesPerUser(user *User, maxFilesPerUser int) {
	if !f.raw.IsValid {
		return
	}
	isValid := false
	message := "User's max files limit reached"

	if (len(user.Files)) < maxFilesPerUser {
		isValid = true
		message = ""
	}

	f.raw.IsValid = isValid
	f.raw.InvalidityMessage = message
}

func (f *File) addRawFile(raw *multipart.FileHeader) {
	f.raw.File = raw
	f.raw.IsValid = true
}
func (f *File) getRawFile() *multipart.FileHeader {
	return f.raw.File
}

func getPrefix(s string) string {
	slices := strings.Split(s, ".")
	prefix := slices[0]
	return prefix
}

func getExtention(s string) string {
	slices := strings.Split(s, ".")
	if len(slices) <= 1 {
		return ""
	}
	format := slices[len(slices)-1]
	return format
}

func generateUUID() string {
	return uuid.New().String()
}
