package database

import (
	"github.com/nu12/audio-gonverter/internal/file"
	"github.com/nu12/audio-gonverter/internal/user"
)

type MockDB struct {
	Messages []string
}

func (*MockDB) Save(*user.User) error {
	return nil
}

func (*MockDB) Load(id string) (*user.User, error) {
	return &user.User{
		UUID:         id,
		IsUploading:  false,
		IsConverting: false,
		Files: []*file.File{
			{
				OriginalName:  "MockFile1.mp3",
				ConvertedName: "MockFile1.ogg",
				OriginalSize:  4011,
				ConvertedSize: 5512,
				OriginalId:    "xxx-xxx-xxx-xxx",
				ConvertedId:   "yyy-yyy-yyy-yyy",
				IsConverted:   true,
			},
			{
				OriginalName: "MockFile1.mp3",
				OriginalSize: 8752,
				OriginalId:   "zzz-zzz-zzz-zzz",
				IsConverted:  false,
			},
		},
	}, nil
}
func (*MockDB) Exist(id string) (bool, error) {
	return true, nil
}
