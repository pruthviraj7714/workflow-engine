package rabbitmq

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Consumer struct {
	MQ *RabbitMQ
}

func NewConsumer(mq *RabbitMQ) *Consumer {
	return &Consumer{
		MQ: mq,
	}
}

func (c *Consumer) Start(ctx context.Context, handler func(context.Context, uuid.UUID) error) error {
	msgs, err := c.MQ.Channel.ConsumeWithContext(
		ctx,
		WorkflowExecutionQueue, // Queue name
		"",                     // Consumer identifier (empty auto-generates one)
		true,                   // Auto-Ack (true automatically acknowledges message receipt)
		false,                  // Exclusive
		false,                  // No-local
		false,                  // No-wait
		nil,                    // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	//Keep consuming messages asynchronously using a channel loop
	forever := make(chan struct{})

	go func() {
		for d := range msgs {

			executionID, err := uuid.Parse(string(d.Body))

			if err != nil {
				log.Printf("invalid UUID: %v", err)
				continue
			}

			if err := handler(ctx, executionID); err != nil {
				log.Printf("handler error: %v", err)
			}

			log.Printf(" [x] successfully executed execution with id : %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
