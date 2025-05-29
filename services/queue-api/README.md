# Queue API Service

## Overview

The Queue API Service provides centralized message queuing capabilities for the microservices architecture. It integrates with shared libraries to provide robust message handling with proper monitoring, logging, and error handling.

## Architecture

### Core Components

1. **Queue Layer**

   - Message operations
   - Message persistence
   - Dead letter handling
   - Connection pooling

2. **Service Layer**

   - Business logic
   - Message transformation
   - Integration with shared libraries
   - Error handling

3. **Integration Layer**
   - Shared libraries integration
   - API services communication
   - Circuit breaking
   - Retry mechanisms

### Shared Libraries Integration

1. **Logging Library**

   ```go
   // Initialize logger
   logger := logging.NewLogger("queue-api")

   // Usage example
   logger.Info("Processing message",
       logging.WithField("queue", "default"),
       logging.WithField("message_id", messageID))
   ```

2. **Monitoring Library**

   ```go
   // Initialize collector
   monitor := monitoring.NewCollector("queue-api")

   // Usage example
   monitor.IncQueueMessages("default")
   defer monitor.ObserveDuration("queue_process")
   ```

3. **Queue Library**

   ```go
   // Initialize queue client
   queueClient := queue.NewAPIClient(queue.Config{
       Endpoint: "http://queue-api:8080",
       Timeout:  time.Second * 5,
   })

   // Usage example
   err = queueClient.Publish(ctx, "default", &queue.Message{
       Type: "default",
       Data: messageData,
   })
   ```

### Service Integration Library

```go
// Initialize service integration
integration := integration.NewServiceIntegration(integration.Config{
    ServiceName: "queue-api",
    Discovery:   "kubernetes",
})

// Register health checks
integration.RegisterHealthCheck("rabbitmq", func() error {
    return rabbitmq.HealthCheck(ctx)
})

// Use circuit breaker
breaker := integration.NewCircuitBreaker("rabbitmq", integration.CircuitBreakerConfig{
    Threshold: 5,
    Timeout:   time.Second * 30,
})

// Use retry mechanism
retry := integration.NewRetry(integration.RetryConfig{
    MaxAttempts: 3,
    Backoff:     time.Second * 2,
})
```

## Complex Integration Patterns

### 1. Message Processing with Retry

```go
func (s *QueueService) ProcessMessage(ctx context.Context, message *Message) error {
    // Create processing context
    procCtx := integration.NewProcessingContext(ctx)
    defer procCtx.Cleanup()

    // Setup retry mechanism
    retry := integration.NewRetry(integration.RetryConfig{
        MaxAttempts: 3,
        Backoff:     time.Second * 2,
    })

    // Process message with retry
    return retry.Execute(procCtx, func(ctx context.Context) error {
        // Process message
        if err := s.processor.Process(ctx, message); err != nil {
            // Move to dead letter queue if max retries exceeded
            if retry.IsMaxAttempts() {
                return s.deadLetterQueue.Publish(ctx, message)
            }
            return err
        }
        return nil
    })
}
```

### 2. Batch Message Publishing

```go
func (s *QueueService) BatchPublish(ctx context.Context, messages []*Message) error {
    // Setup batch processor
    processor := integration.NewBatchProcessor(integration.BatchProcessorConfig{
        BatchSize: 100,
        Workers:   5,
    })

    // Process messages in batches
    return processor.Process(ctx, messages, func(ctx context.Context, batch []*Message) error {
        return s.queue.BatchPublish(ctx, batch)
    })
}
```

## Configuration

### Base Configuration

```yaml
service:
  name: queue-api
  version: 1.0.0
  port: 8080

logging:
  level: info
  format: json
  output: stdout

monitoring:
  enabled: true
  prometheus:
    path: /metrics
    port: 9090

integration:
  service_discovery: kubernetes
  circuit_breaker:
    threshold: 5
    timeout: 30s
  retry:
    max_attempts: 3
    backoff: 2s
```

### Service-Specific Configuration

```yaml
queue:
  type: rabbitmq
  host: rabbitmq
  port: 5672
  user: queue_api
  password: ${RABBITMQ_PASSWORD}
  max_connections: 20
  connection_timeout: 5s
  circuit_breaker:
    threshold: 5
    timeout: 30s
```

## Error Handling

### Standard Error Patterns

```go
// Error handling with logging and monitoring
if err != nil {
    switch {
    case errors.Is(err, queue.ErrQueueFull):
        logger.Warn("Queue full", logging.WithError(err))
        monitor.IncQueueFull()
    case errors.Is(err, queue.ErrConnection):
        logger.Error("Queue connection error", logging.WithError(err))
        monitor.IncConnectionErrors()
    case errors.Is(err, queue.ErrTimeout):
        logger.Error("Queue timeout", logging.WithError(err))
        monitor.IncTimeoutErrors()
    default:
        logger.Error("Unexpected error", logging.WithError(err))
        monitor.IncErrors()
    }
    return nil, err
}
```

## Health Checks

### Service Health

```go
// Register health checks
integration.RegisterHealthCheck("rabbitmq", func() error {
    return rabbitmq.HealthCheck(ctx)
})
```

## Metrics

### Standard Metrics

```go
// Queue metrics
monitor.IncQueueMessages("default")
monitor.ObserveQueueLatency("default", latency)

// RabbitMQ metrics
monitor.IncRabbitMQOperations("publish")
monitor.ObserveRabbitMQLatency("publish", latency)
monitor.ObserveConnectionPoolSize(poolSize)
```

## Development

### Setup

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Run tests:

   ```bash
   go test ./...
   ```

3. Build service:
   ```bash
   go build -o queue-api ./cmd/queue-api
   ```

### Testing

1. Unit tests:

   ```bash
   go test -v ./internal/...
   ```

2. Integration tests:

   ```bash
   go test -v ./tests/integration/...
   ```

3. Load tests:
   ```bash
   k6 run ./tests/load/queue-api.js
   ```

## Deployment

### Kubernetes

1. Apply configurations:

   ```bash
   kubectl apply -f k8s/
   ```

2. Verify deployment:

   ```bash
   kubectl get pods -n microservices
   ```

3. Check logs:
   ```bash
   kubectl logs -n microservices -l app=queue-api
   ```

### Docker

1. Build image:

   ```bash
   docker build -t queue-api:latest .
   ```

2. Run container:
   ```bash
   docker run -p 8080:8080 queue-api:latest
   ```

## Monitoring

### Prometheus Metrics

- Queue message rates
- Error rates
- Latency percentiles
- Connection pool metrics
- Dead letter queue metrics

### Grafana Dashboards

- Service overview
- Error rates
- Latency trends
- Queue performance
- Dead letter queue status

## Logging

### Log Levels

- ERROR: Service errors
- WARN: Queue full, retries
- INFO: Message processing
- DEBUG: Detailed operations

### Log Fields

- service: queue-api
- trace_id: Request tracing
- message_id: Message identifier
- queue: Queue name
- operation: Operation type
- duration: Operation time
- error: Error details

## Security

### Authentication

- JWT token validation
- Service-to-service authentication
- API key management

### Authorization

- Role-based access control
- Permission management
- Resource access control

## Dependencies

### External Services

- RabbitMQ
- Queue API Service

### Shared Libraries

- Logging Library
- Monitoring Library
- Queue Client Library
- Service Integration Library

## API Documentation

### OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: Queue API
  version: 1.0.0
paths:
  /queue:
    post:
      summary: Publish message
      responses:
        "201":
          description: Created
  /queue/{queue}:
    get:
      summary: Consume message
      responses:
        "200":
          description: Success
    post:
      summary: Publish message to queue
      responses:
        "201":
          description: Created
  /queue/{queue}/dead-letter:
    get:
      summary: Get dead letter messages
      responses:
        "200":
          description: Success
    post:
      summary: Republish dead letter message
      responses:
        "201":
          description: Created
```

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create pull request

## License

MIT License

## Implementation Status

### Current State

1. **Queue Layer**

   - [ ] Queue client implementation
   - [ ] Queue configuration
   - [ ] Message type definitions
   - [ ] Error handling
   - [ ] Metrics collection

2. **Service Layer**

   - [ ] Business logic implementation
   - [ ] Message transformation
   - [ ] Shared libraries integration
   - [ ] Error handling

3. **Integration Layer**
   - [ ] Shared libraries integration
   - [ ] API services communication
   - [ ] Circuit breaking
   - [ ] Retry mechanisms

### Implementation Plan

1. **Phase 1: Core Infrastructure**

   - [ ] Project structure setup
   - [ ] Configuration management
   - [ ] Logging integration
   - [ ] Metrics collection

2. **Phase 2: Queue Implementation**

   - [ ] RabbitMQ client setup
   - [ ] Queue operations
   - [ ] Error handling
   - [ ] Metrics collection

3. **Phase 3: Service Integration**
   - [ ] Shared libraries integration
   - [ ] API services communication
   - [ ] Circuit breaking
   - [ ] Retry mechanisms

## API Endpoints

### 1. Queue Operations

```http
POST /api/v1/queue/{queue_name}
GET /api/v1/queue/{queue_name}
DELETE /api/v1/queue/{queue_name}/{message_id}
```

### 2. Queue Management

```http
GET /api/v1/queue/stats
POST /api/v1/queue/purge
GET /api/v1/queue/status
```

### 3. Health and Metrics

```http
GET /health
GET /ready
GET /metrics
```

## Error Types

### 1. Queue Errors

- Connection errors
- Operation errors
- Message errors
- Timeout errors
- Serialization errors

### 2. Service Errors

- Validation errors
- Dependency errors
- Resource errors
- Timeout errors

### Recovery Strategies

### 1. Queue Recovery

- Connection retry
- Operation retry
- Message retry
- Error logging
- Circuit breaking

### 2. Service Recovery

- Validation recovery
- Dependency recovery
- Resource recovery
- Error logging
- Circuit breaking

## Cross-References

- [Queue Service Patterns](../../reference-materials/development/patterns/queue-service-patterns.md)
- [Service Integration Patterns](../../reference-materials/development/patterns/service-integration-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Error Handling Patterns](../../reference-materials/development/patterns/error-handling-patterns.md)
