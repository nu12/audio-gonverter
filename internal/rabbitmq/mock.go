package rabbitmq

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
