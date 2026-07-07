package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func StartProducer() {
	// Configure the writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"), // Default Kafka port
		Topic:        "user-events",
		Balancer:     &kafka.LeastBytes{}, // Automatically balances messages across partitions
		BatchTimeout: 10 * time.Millisecond,
	}
	defer writer.Close()

	// Publish a message

	// if err != nil {
	// 	log.Fatalf("Could not write message: %v", err)
	// }

	return writer
}
