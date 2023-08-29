package model

import "strings"

type File struct {
	OriginalName  string
	ConvertedName string
	OriginalSize  int64
	ConvertedSize int64
	OriginalId    string
	ConvertedId   string
	IsConverted   bool
}

func NewFile(OriginalName string) (*File, error) {
	slices := strings.Split(OriginalName, ".")
	format := slices[len(slices)-1]

	return &File{
		OriginalName: OriginalName,
		OriginalId:   "uuid." + format,
	}, nil
}
