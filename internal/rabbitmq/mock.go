package rabbitmq

type QueueMock struct{}

func (*QueueMock) Push(msg string) error {
	return nil
}

func (*QueueMock) Consume() <-chan any {
	return nil
}
