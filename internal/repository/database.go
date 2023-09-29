package repository

import (
	"github.com/nu12/audio-gonverter/internal/user"
)

type DatabaseRepository interface {
	Save(*user.User) error
	Load(string) (*user.User, error)
	Exist(string) (bool, error)
}
