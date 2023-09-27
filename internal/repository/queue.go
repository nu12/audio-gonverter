package repository

type QueueMessage struct {
	UserUUID string `json:"user"`
	Format   string `json:"format"`
	Kbps     string `json:"kbps"`
}

type QueueRepository interface {
	Push(string) error
	Pull() (string, error)
	Encode(QueueMessage) (string, error)
	Decode(string) (QueueMessage, error)
}
