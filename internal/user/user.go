package user

import (
	"github.com/google/uuid"
	"github.com/nu12/audio-gonverter/internal/file"
)

type User struct {
	UUID         string       `json:"uuid"`
	Files        []*file.File `json:"files"`
	IsUploading  bool         `json:"is_uploading"`
	IsConverting bool         `json:"is_converting"`
	Messages     []string     `json:"messages"`
}

func NewUser() User {
	return User{
		UUID:         GenerateUUID(),
		Files:        []*file.File{},
		IsUploading:  false,
		IsConverting: false,
	}
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
