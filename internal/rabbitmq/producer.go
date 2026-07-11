package rabbitmq

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	MQ *RabbitMQ
}

func NewProducer(mq *RabbitMQ) *Producer {
	return &Producer{
		MQ: mq,
	}
}

func (p *Producer) PublishWorkflowExecution(executionID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := executionID.String()

	err := p.MQ.Channel.PublishWithContext(ctx,
		"",                     // Exchange name (empty string uses default direct exchange)
		WorkflowExecutionQueue, // Routing key (queue name)
		false,                  // Mandatory
		false,                  // Immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent, // Make message persistent
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}
