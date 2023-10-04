package event

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("pushing to channel")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"logs_topic", // exchange
		severity,     // key
		false,        // mandatory
		false,        // immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	return err
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, nil
	}
	return emitter, nil
}
