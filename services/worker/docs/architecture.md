# Worker Service Architecture

## Overview

The Worker Service is designed to handle asynchronous task processing for the Profile Service, specifically managing email validation and image generation tasks. This document outlines the architectural decisions, patterns, and implementation details.

## Architecture Diagram

```mermaid
graph TD
    A[Profile API] -->|Publish| B[RabbitMQ]
    B -->|Consume| C[Worker Service]
    C -->|Process| D[Email Processor]
    C -->|Process| E[Image Processor]
    D -->|Update| F[Profile Storage]
    E -->|Update| F
    D -->|Send| G[Email Service]
    E -->|Generate| H[AI Service]
    C -->|Metrics| I[Monitoring]
    C -->|Health| J[Health Checks]
```

## Core Components

### 1. Message Queue Integration

```mermaid
sequenceDiagram
    participant API as Profile API
    participant Queue as RabbitMQ
    participant Worker as Worker Service
    participant Storage as Profile Storage

    API->>Queue: Publish Task
    Queue->>Worker: Consume Message
    Worker->>Worker: Process Task
    Worker->>Storage: Update Profile
    Storage-->>Worker: Confirm Update
    Worker-->>Queue: Acknowledge
```

#### Queue Structure

- **Email Queue**

  - Name: `profile-email-queue`
  - Purpose: Email validation tasks
  - Durability: Yes
  - TTL: 24 hours
  - Dead Letter Exchange: Yes

- **Image Queue**
  - Name: `profile-image-queue`
  - Purpose: Image generation tasks
  - Durability: Yes
  - TTL: 24 hours
  - Dead Letter Exchange: Yes

### 2. Task Processing

#### Email Validation Flow

```mermaid
sequenceDiagram
    participant Queue as RabbitMQ
    participant Worker as Worker Service
    participant Email as Email Service
    participant Storage as Profile Storage

    Queue->>Worker: Email Task
    Worker->>Email: Send Validation
    Email-->>Worker: Email Sent
    Worker->>Storage: Update Status
    Storage-->>Worker: Confirmed
    Worker-->>Queue: Acknowledge
```

#### Image Generation Flow

```mermaid
sequenceDiagram
    participant Queue as RabbitMQ
    participant Worker as Worker Service
    participant AI as AI Service
    participant Storage as Profile Storage

    Queue->>Worker: Image Task
    Worker->>AI: Generate Image
    AI-->>Worker: Image URL
    Worker->>Storage: Update Profile
    Storage-->>Worker: Confirmed
    Worker-->>Queue: Acknowledge
```

### 3. Error Handling

#### Error Recovery Flow

```mermaid
graph TD
    A[Error Occurs] --> B{Recoverable?}
    B -->|Yes| C[Retry Task]
    B -->|No| D[Dead Letter]
    C --> E{Max Retries?}
    E -->|No| F[Process Again]
    E -->|Yes| D
    D --> G[Alert]
    D --> H[Log]
```

### 4. Monitoring

#### Metrics Collection

```mermaid
graph TD
    A[Worker Service] --> B[Queue Metrics]
    A --> C[Processing Metrics]
    A --> D[System Metrics]
    B --> E[Prometheus]
    C --> E
    D --> E
    E --> F[Grafana]
    E --> G[Alerting]
```

## Implementation Details

### 1. Service Structure

```
worker-service/
├── cmd/
│   └── worker/
│       └── main.go
├── internal/
│   ├── queue/
│   │   ├── consumer.go
│   │   ├── publisher.go
│   │   └── config.go
│   ├── processor/
│   │   ├── email.go
│   │   └── image.go
│   ├── storage/
│   │   └── client.go
│   ├── monitoring/
│   │   ├── metrics.go
│   │   └── health.go
│   └── config/
│       └── config.go
└── pkg/
    └── shared/
        ├── queue.go
        └── monitoring.go
```

### 2. Key Interfaces

```go
// Queue Consumer Interface
type Consumer interface {
    Consume(ctx context.Context) error
    Close() error
}

// Task Processor Interface
type Processor interface {
    Process(ctx context.Context, task Task) error
    Validate(task Task) error
}

// Storage Client Interface
type StorageClient interface {
    UpdateProfile(ctx context.Context, profile Profile) error
    GetProfile(ctx context.Context, id string) (Profile, error)
}
```

### 3. Configuration Management

```yaml
worker:
  queue:
    host: rabbitmq
    port: 5672
    username: worker
    password: ${QUEUE_PASSWORD}
    vhost: /workers
    prefetch: 10
    reconnect:
      maxAttempts: 5
      initialDelay: 1s
      maxDelay: 30s

  processing:
    maxRetries: 3
    retryDelay: 5s
    timeout: 30s
    batchSize: 10
    concurrency: 5

  monitoring:
    metricsPort: 9090
    healthCheckInterval: 30s
    logLevel: info
```

## Cross-References

- [Worker Service Patterns](../../reference-materials/development/patterns/worker-service-patterns.md)
- [Queuing Patterns](../../reference-materials/development/patterns/queuing-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Long-Running Tasks](../../reference-materials/development/patterns/long-running-tasks.md)
