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
		UUID:         generateUUID(),
		Files:        []File{},
		IsUploading:  false,
		IsConverting: false,
	}
}

func (u *User) AddFile(file *File) error {
	files := u.Files
	files = append(files, *file)
	u.Files = files
	return nil
}

func (u *User) RemoveFile(id string) error {
	files := []File{}
	for _, f := range u.Files {
		if f.OriginalId != id {
			files = append(files, f)
		}
	}
	u.Files = files
	return nil
}

func (u *User) ClearFiles() error {
	u.Files = []File{}
	return nil
}
