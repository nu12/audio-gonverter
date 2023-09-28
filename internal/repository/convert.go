package repository

import "github.com/nu12/audio-gonverter/internal/model"

type ConvertionToolRepo interface {
	Convert(file *model.File, format, kpbs string) error
}
