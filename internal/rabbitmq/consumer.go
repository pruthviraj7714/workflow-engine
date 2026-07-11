package rabbitmq

import (
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

func (c *Consumer) Start(handler func(uuid.UUID) error) error {
	msgs, err := c.MQ.Channel.Consume(
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
			log.Printf(" [x] Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
