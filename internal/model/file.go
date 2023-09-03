package model

import (
	"strings"

	"github.com/google/uuid"
)

type File struct {
	OriginalName  string `json:"original_name"`
	ConvertedName string `json:"converted_name"`
	OriginalSize  int64  `json:"original_size"`
	ConvertedSize int64  `json:"converted_size"`
	OriginalId    string `json:"original_id"`
	ConvertedId   string `json:"converted_id"`
	IsConverted   bool   `json:"is_converted"`
}

func NewFile(OriginalName string) (*File, error) {

	return &File{
		OriginalName: OriginalName,
		OriginalId:   generateUUID() + "." + getExtention(OriginalName),
	}, nil
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
