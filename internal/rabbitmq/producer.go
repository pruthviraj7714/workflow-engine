package rabbitmq

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a unique channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare a queue to send to
	q, err := ch.QueueDeclare(
		"task_queue", // Queue name
		true,         // Durable (survives server restarts)
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 4. Define context timeout for publishing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World from Go!"

	// 5. Publish the message to the default exchange
	err = ch.PublishWithContext(ctx,
		"",     // Exchange name (empty string uses default direct exchange)
		q.Name, // Routing key (queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s", body)
}
