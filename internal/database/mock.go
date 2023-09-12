package database

import "github.com/nu12/audio-gonverter/internal/model"

type MockDB struct{}

func (*MockDB) Save(*model.User) error {
	return nil
}

func (*MockDB) Load(id string) (*model.User, error) {
	return &model.User{
		UUID:         id,
		IsUploading:  false,
		IsConverting: false,
		Files: []model.File{
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
