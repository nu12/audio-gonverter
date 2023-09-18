package repository

type QueueRepository interface {
	Push(string) error
	Pull() (string, error)
}
