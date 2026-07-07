package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Configure the reader as part of a consumer group
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "email-service-group", // Consumer group ID
		Topic:    "user-events",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	log.Println("Consumer started, listening for messages...")

	// Continuous loop to fetch events
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Process the event
		log.Printf("Received message: Key=%s, Value=%s, Partition=%d, Offset=%d\n",
			string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
	}
}
