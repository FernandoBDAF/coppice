package main

import (
	"log"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/base"
	"github.com/fernandobarroso/microservices/operational-workers/internal/common/utils"
	"github.com/fernandobarroso/microservices/operational-workers/internal/processors/profile"
)

func main() {
	log.Println("Starting Profile Worker...")

	// profile.task: queue/routing key per CONTRACTS.md §2. The EXCHANGE is
	// intentionally "tasks-exchange", NOT "profile-tasks" as CONTRACTS.md /
	// ROUTING_KEYS.md state.
	//
	// api-service's actual publisher config
	// (api-service/internal/domain/task/model.go DefaultRoutingMap["profile.task"])
	// declares Exchange: "tasks-exchange", TTL: 24h, DeadLetterTTL: 7*24h,
	// MaxRetries: 3 — not "profile-tasks" / 1h as the docs claim. RabbitMQ
	// requires the consumer's queue-declare args to be byte-identical to the
	// publisher's, and the queue only ever receives messages published to the
	// exchange it is bound to, so this worker matches the PUBLISHER's real
	// behavior. Flagged for the orchestrator: either rename api-service's
	// exchange to "profile-tasks" (and fix its TTL) or update CONTRACTS.md /
	// ROUTING_KEYS.md to document "tasks-exchange" / 24h.
	config := &base.WorkerConfig{
		WorkerType:    "profile",
		QueueName:     "profile-processing",
		ExchangeName:  "tasks-exchange",
		RoutingKey:    "profile.task",
		PrefetchCount: 2,
		MessageTTL:    24 * time.Hour,
		DeadLetterTTL: 7 * 24 * time.Hour,
		MaxRetries:    3,
		HTTPPort:      utils.GetEnvOrDefault("HEALTH_PORT", "8080"),
	}

	processor := profile.NewProcessor()

	worker, err := base.NewBaseWorker(config, processor)
	if err != nil {
		log.Fatalf("Failed to create profile worker: %v", err)
	}

	if err := worker.Run(); err != nil {
		log.Fatalf("Profile worker failed: %v", err)
	}
}
