# Worker Service API Documentation

## Overview

This document outlines the API interfaces and message formats used by the Worker Service for task processing and service communication.

## Message Types

### 1. Email Validation Message

```protobuf
message EmailValidationMessage {
    string profile_id = 1;
    string email = 2;
    string validation_token = 3;
    int64 created_at = 4;
    map<string, string> metadata = 5;
}
```

#### Fields

| Field            | Type                | Description                        |
| ---------------- | ------------------- | ---------------------------------- |
| profile_id       | string              | Unique identifier of the profile   |
| email            | string              | Email address to validate          |
| validation_token | string              | Token for email validation         |
| created_at       | int64               | Unix timestamp of message creation |
| metadata         | map<string, string> | Additional message metadata        |

### 2. Image Generation Message

```protobuf
message ImageGenerationMessage {
    string profile_id = 1;
    string prompt = 2;
    string style = 3;
    int64 created_at = 4;
    map<string, string> metadata = 5;
}
```

#### Fields

| Field      | Type                | Description                        |
| ---------- | ------------------- | ---------------------------------- |
| profile_id | string              | Unique identifier of the profile   |
| prompt     | string              | Image generation prompt            |
| style      | string              | Image style specification          |
| created_at | int64               | Unix timestamp of message creation |
| metadata   | map<string, string> | Additional message metadata        |

## Queue Configuration

### 1. Email Queue

```yaml
queues:
  email:
    name: profile-email-queue
    exchange: profile-events
    routing_key: profile.email.*
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 86400000 # 24 hours
      x-dead-letter-exchange: profile-dlx
      x-dead-letter-routing-key: profile.email.dlq
```

### 2. Image Queue

```yaml
queues:
  image:
    name: profile-image-queue
    exchange: profile-events
    routing_key: profile.image.*
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 86400000 # 24 hours
      x-dead-letter-exchange: profile-dlx
      x-dead-letter-routing-key: profile.image.dlq
```

## Service Interfaces

### 1. Queue Consumer Interface

```go
type Consumer interface {
    // Consume starts consuming messages from the queue
    Consume(ctx context.Context) error

    // Close gracefully shuts down the consumer
    Close() error

    // GetQueueName returns the name of the queue
    GetQueueName() string

    // GetMessageCount returns the current message count
    GetMessageCount() (int, error)
}
```

### 2. Task Processor Interface

```go
type Processor interface {
    // Process handles the task processing
    Process(ctx context.Context, task Task) error

    // Validate validates the task before processing
    Validate(task Task) error

    // GetTaskType returns the type of task this processor handles
    GetTaskType() string

    // GetMetrics returns processor-specific metrics
    GetMetrics() map[string]float64
}
```

### 3. Storage Client Interface

```go
type StorageClient interface {
    // UpdateProfile updates the profile with new data
    UpdateProfile(ctx context.Context, profile Profile) error

    // GetProfile retrieves a profile by ID
    GetProfile(ctx context.Context, id string) (Profile, error)

    // UpdateProfileStatus updates only the profile status
    UpdateProfileStatus(ctx context.Context, id string, status string) error
}
```

## Error Types

### 1. Queue Errors

```go
type QueueError struct {
    Code    string
    Message string
    Cause   error
}

const (
    ErrConnectionFailed = "CONNECTION_FAILED"
    ErrChannelFailed    = "CHANNEL_FAILED"
    ErrMessageInvalid   = "MESSAGE_INVALID"
    ErrConsumerFailed   = "CONSUMER_FAILED"
)
```

### 2. Processing Errors

```go
type ProcessingError struct {
    Code    string
    Message string
    TaskID  string
    Cause   error
}

const (
    ErrValidationFailed = "VALIDATION_FAILED"
    ErrProcessingFailed = "PROCESSING_FAILED"
    ErrStorageFailed    = "STORAGE_FAILED"
    ErrIntegrationFailed = "INTEGRATION_FAILED"
)
```

## Metrics

### 1. Queue Metrics

```go
type QueueMetrics struct {
    MessagesProcessed    prometheus.Counter
    MessagesFailed      prometheus.Counter
    ProcessingDuration  prometheus.Histogram
    QueueDepth         prometheus.Gauge
    ConsumerLag        prometheus.Gauge
}
```

### 2. Processing Metrics

```go
type ProcessingMetrics struct {
    TasksProcessed     prometheus.Counter
    TasksFailed       prometheus.Counter
    ProcessingTime    prometheus.Histogram
    RetryCount        prometheus.Counter
    ErrorTypes        *prometheus.CounterVec
}
```

## Health Checks

### 1. Queue Health

```go
type QueueHealth struct {
    ConnectionStatus string
    ChannelStatus   string
    QueueStatus     string
    ConsumerStatus  string
    LastError       error
}
```

### 2. Service Health

```go
type ServiceHealth struct {
    Status          string
    QueueHealth     QueueHealth
    StorageHealth   StorageHealth
    ProcessorHealth ProcessorHealth
    LastChecked     time.Time
}
```

## Cross-References

- [Worker Service Patterns](../../reference-materials/development/patterns/worker-service-patterns.md)
- [Queuing Patterns](../../reference-materials/development/patterns/queuing-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Long-Running Tasks](../../reference-materials/development/patterns/long-running-tasks.md)
