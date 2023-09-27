package rabbitmq

import (
	"github.com/nu12/audio-gonverter/internal/repository"
)

type QueueMock struct {
	Count int
}

func (q *QueueMock) Push(msg string) error {
	q.Count++
	return nil
}

func (*QueueMock) Pull() (string, error) {
	return "Queue mock", nil
}

func (*QueueMock) Encode(msg repository.QueueMessage) (string, error) {
	return "", nil
}

func (*QueueMock) Decode(msg string) (repository.QueueMessage, error) {
	return repository.QueueMessage{}, nil
}
