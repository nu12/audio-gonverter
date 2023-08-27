package model

type File struct {
	OriginalName  string
	ConvertedName string
	OriginalSize  int64
	ConvertedSize int64
	OriginalId    string
	ConvertedId   string
	IsConverted   bool
}

func NewFile(OriginalName string, OriginalSize int64) (*File, error) {
	// Create random ID
	return &File{}, nil
}
