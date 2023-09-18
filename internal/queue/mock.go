package queue

type QueueMock struct{}

func (*QueueMock) Push(msg string) error {
	return nil
}

func (*QueueMock) Pull() (string, error) {
	return "", nil
}
