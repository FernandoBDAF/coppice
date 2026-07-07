package main

import (
	"log"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/base"
	"github.com/fernandobarroso/microservices/operational-workers/internal/common/utils"
	"github.com/fernandobarroso/microservices/operational-workers/internal/processors/image"
)

func main() {
	log.Println("Starting Image Worker...")

	// image.process: exchange/queue/routing key per CONTRACTS.md §2 and
	// graph-worker/shared/contracts/ROUTING_KEYS.md. TTL/DeadLetterTTL/
	// MaxRetries mirror api-service's DefaultRoutingMap["image.process"]
	// (internal/domain/task/model.go) exactly, since RabbitMQ requires
	// identical queue-declare args between publisher and consumer.
	config := &base.WorkerConfig{
		WorkerType:    "image",
		QueueName:     "image-processing",
		ExchangeName:  "image-tasks",
		RoutingKey:    "image.process",
		PrefetchCount: 1, // resource intensive - process one at a time
		MessageTTL:    6 * time.Hour,
		DeadLetterTTL: 3 * 24 * time.Hour,
		MaxRetries:    2,
		HTTPPort:      utils.GetEnvOrDefault("HEALTH_PORT", "8080"),
	}

	processor := image.NewImageProcessor()

	worker, err := base.NewBaseWorker(config, processor)
	if err != nil {
		log.Fatalf("Failed to create image worker: %v", err)
	}

	if err := worker.Run(); err != nil {
		log.Fatalf("Image worker failed: %v", err)
	}
}
