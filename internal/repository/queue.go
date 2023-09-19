package repository

type QueueRepository interface {
	Push(string) error
}
