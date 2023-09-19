package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nu12/audio-gonverter/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	User   *model.User `json:"user"`
	Format string      `json:"format"`
	Kbps   string      `json:"kbps"`
}
type RabbitQueue struct {
	Channel  *amqp.Channel
	Queue    amqp.Queue
	Consumer <-chan amqp.Delivery
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

func (q *RabbitQueue) Consume() <-chan amqp.Delivery {
	return q.Consumer
}

func Connect(connString string) (*RabbitQueue, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return &RabbitQueue{}, err
	}
	//defer ch.Close()

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
		Channel:  ch,
		Queue:    q,
		Consumer: c,
	}, nil

}

func Encode(msg Message) (string, error) {
	j, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func Decore(msg string) (Message, error) {
	return Message{}, nil
}
