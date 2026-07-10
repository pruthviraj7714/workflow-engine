package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func StartConsumer() {
	// 1. Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare the queue (idempotent; ensures it exists if consumer starts first)
	q, err := ch.QueueDeclare(
		"task_queue", // Queue name
		true,         // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 4. Register as a consumer
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer identifier (empty auto-generates one)
		true,   // Auto-Ack (true automatically acknowledges message receipt)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 5. Keep consuming messages asynchronously using a channel loop
	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf(" [x] Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
