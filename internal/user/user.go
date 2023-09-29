package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nu12/audio-gonverter/internal/file"
)

type User struct {
	UUID         string       `json:"uuid"`
	Files        []*file.File `json:"files"`
	IsUploading  bool         `json:"is_uploading"`
	IsConverting bool         `json:"is_converting"`
	Messages     []string     `json:"messages"`

	err error `json:"-"`
}

func New() *User {
	return &User{
		UUID:         GenerateUUID(),
		Files:        []*file.File{},
		IsUploading:  false,
		IsConverting: false,
		err:          nil,
	}
}

// SA1029: Users of WithValue should define their own types for keys.
type userID string

func FromRequest(r *http.Request) *User {
	u := r.Context().Value(userID("user"))
	if u == nil {
		return &User{err: errors.New("User not found in Request")}
	}
	return u.(*User)
}

func (u *User) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userID("user"), u)
}

func (u *User) Err() error {
	return u.err
}

func (u *User) AddFile(file *file.File) error {
	files := u.Files
	files = append(files, file)
	u.Files = files
	return nil
}

func (u *User) RemoveFile(id string) error {
	files := []*file.File{}
	for _, f := range u.Files {
		if f.OriginalId != id {
			files = append(files, f)
		}
	}
	u.Files = files
	return nil
}

func (u *User) ClearFiles() error {
	u.Files = []*file.File{}
	return nil
}

func (u *User) AddMessage(s string) {
	u.Messages = append(u.Messages, s)
}

func (u *User) GetMessages() []string {
	if len(u.Messages) == 0 {
		return []string{"Welcome to audio-gonverter!"}
	}
	messages := u.Messages
	u.Messages = []string{}
	return messages
}

func GenerateUUID() string {
	return uuid.New().String()
}
