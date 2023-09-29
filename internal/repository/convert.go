package repository

import "github.com/nu12/audio-gonverter/internal/file"

type ConvertionToolRepo interface {
	Convert(file *file.File, format, kpbs string) error
}
