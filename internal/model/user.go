package model

type User struct {
	UUID         string
	Files        []File
	IsUploading  bool
	IsConverting bool
	// Add a DB Repo here
}

func NewUser() User {
	return User{
		UUID: "to-be-implemented",
		Files: []File{
			{OriginalName: "Test1"},
			{OriginalName: "Test2"},
		},
		IsUploading:  false,
		IsConverting: false,
	}

	//save user
}

func (u *User) AddFile(file *File) error {
	return nil
}

func (u *User) RemoveFile(id string) error {
	return nil
}

func (u *User) ClearFiles(id string) error {
	u.Files = []File{}
	return nil
}
