package main

import (
	"log"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/base"
	"github.com/fernandobarroso/microservices/operational-workers/internal/common/utils"
	"github.com/fernandobarroso/microservices/operational-workers/internal/processors/email"
)

func main() {
	log.Println("Starting Email Worker...")

	// email.send: exchange/queue/routing key per CONTRACTS.md §2 and
	// graph-worker/shared/contracts/ROUTING_KEYS.md. TTL/DeadLetterTTL/
	// MaxRetries mirror api-service's DefaultRoutingMap["email.send"]
	// (internal/domain/task/model.go) exactly, since RabbitMQ requires
	// identical queue-declare args between publisher and consumer.
	config := &base.WorkerConfig{
		WorkerType:    "email",
		QueueName:     "email-processing",
		ExchangeName:  "email-tasks",
		RoutingKey:    "email.send",
		PrefetchCount: 5, // higher throughput for burst processing
		MessageTTL:    1 * time.Hour,
		DeadLetterTTL: 24 * time.Hour,
		MaxRetries:    5,
		HTTPPort:      utils.GetEnvOrDefault("HEALTH_PORT", "8080"),
	}

	processor := email.NewEmailProcessor()

	worker, err := base.NewBaseWorker(config, processor)
	if err != nil {
		log.Fatalf("Failed to create email worker: %v", err)
	}

	if err := worker.Run(); err != nil {
		log.Fatalf("Email worker failed: %v", err)
	}
}
