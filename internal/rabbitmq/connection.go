package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		WorkflowExecutionQueue, // Queue name
		true,                   // Durable
		false,                  // Delete when unused
		false,                  // Exclusive
		false,                  // No-wait
		nil,                    // Arguments
	)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	defer func() {
		r.Conn.Close()

		r.Channel.Close()
	}()
}
