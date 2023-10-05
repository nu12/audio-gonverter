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

func (u *User) AddFile(file *file.File) *User {
	if u.err != nil {
		return u
	}

	files := u.Files
	files = append(files, file)
	u.Files = files
	return u
}

func (u *User) RemoveFile(id string) *User {
	if u.err != nil {
		return u
	}

	files := []*file.File{}
	for _, f := range u.Files {
		if f.OriginalId != id {
			files = append(files, f)
		}
	}
	u.Files = files
	return u
}

func (u *User) ClearFiles() *User {
	if u.err != nil {
		return u
	}

	u.Files = []*file.File{}
	return u
}

func (u *User) AddMessage(s string) *User {
	if u.err != nil {
		return u
	}

	u.Messages = append(u.Messages, s)
	return u
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
