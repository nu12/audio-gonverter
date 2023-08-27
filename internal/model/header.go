package model

import "mime/multipart"

type RawFile interface {
	Open() (multipart.File, error)
	Filename() string
}

type Header struct {
	FileHeader *multipart.FileHeader
}

func (h *Header) Filename() string {
	return h.FileHeader.Filename
}

func (h *Header) Open() (multipart.File, error) {
	return h.FileHeader.Open()
}
