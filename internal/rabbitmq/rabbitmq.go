package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	UserUUID string `json:"user"`
	Format   string `json:"format"`
	Kbps     string `json:"kbps"`
}
type RabbitQueue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Consumer   <-chan amqp.Delivery
}

func (q *RabbitQueue) Push(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return q.Channel.PublishWithContext(ctx,
		"",           // exchange
		q.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

}

func (q *RabbitQueue) Pull() (string, error) {
	for msg := range q.Consumer {
		return string(msg.Body), nil
	}
	return "", errors.New("No message found")
}

func Connect(connString string) (*RabbitQueue, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return &RabbitQueue{}, err
	}

	q, err := ch.QueueDeclare(
		"audio", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments)
	)
	if err != nil {
		return &RabbitQueue{}, err
	}

	c, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return &RabbitQueue{}, err
	}

	return &RabbitQueue{
		Connection: conn,
		Channel:    ch,
		Queue:      q,
		Consumer:   c,
	}, nil

}

func Encode(msg Message) (string, error) {
	j, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func Decode(msg string) (Message, error) {
	var message Message
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}
