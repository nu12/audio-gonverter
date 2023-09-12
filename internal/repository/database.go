package repository

import "github.com/nu12/audio-gonverter/internal/model"

type DatabaseRepository interface {
	Save(*model.User) error
	Load(string) (*model.User, error)
}
