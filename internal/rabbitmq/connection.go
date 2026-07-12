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

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			return err
		}
	}

	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
