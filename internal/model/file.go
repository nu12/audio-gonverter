package model

import (
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

type File struct {
	OriginalName   string                    `json:"original_name"`
	ConvertedName  string                    `json:"converted_name"`
	OriginalSize   int64                     `json:"original_size"`
	ConvertedSize  int64                     `json:"converted_size"`
	OriginalId     string                    `json:"original_id"`
	ConvertedId    string                    `json:"converted_id"`
	IsConverted    bool                      `json:"is_converted"`
	SaveToDiskFunc func(*File, string) error `json:"-"`
	raw            *multipart.FileHeader     `json:"-"`
}

func NewFile(OriginalName string) (*File, error) {
	file := &File{
		OriginalName: OriginalName,
		OriginalId:   generateUUID() + "." + getExtention(OriginalName),
	}
	file.SaveToDiskFunc = func(f *File, path string) error {
		mf, err := f.getRawFile().Open()
		if err != nil {
			return err
		}
		defer mf.Close()

		of, err := os.OpenFile(path+"/"+f.OriginalId, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer of.Close()

		bytes, err := io.Copy(of, mf)
		if err != nil {
			return err
		}

		f.OriginalSize = bytes
		return nil
	}

	return file, nil
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

func (f *File) addRawFile(raw *multipart.FileHeader) {
	f.raw = raw
}
func (f *File) getRawFile() *multipart.FileHeader {
	if f.raw != nil {
		return f.raw
	}
	return &multipart.FileHeader{}
}

func getPrefix(s string) string {
	slices := strings.Split(s, ".")
	prefix := slices[0]
	return prefix
}

func getExtention(s string) string {
	slices := strings.Split(s, ".")
	format := slices[len(slices)-1]
	return format
}

func generateUUID() string {
	return uuid.New().String()
}
