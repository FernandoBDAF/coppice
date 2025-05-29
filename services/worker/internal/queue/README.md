# Profile Queue Service

## Overview

The Profile Queue Service is a message broker service that manages communication between the Profile API and worker services (email and image generation). It provides reliable message delivery, proper routing, and message persistence using RabbitMQ.

## Architecture

### Components

1. **Queue Manager**

   - Manages RabbitMQ connections and channels
   - Handles queue and exchange setup
   - Implements connection pooling
   - Provides health monitoring

2. **Message Processor**

   - Validates incoming messages
   - Transforms messages as needed
   - Routes messages to appropriate queues
   - Handles message persistence
   - Manages message priorities
   - Implements message correlation
   - Provides message tracing
   - Supports message replay

3. **Service Integrator**

   - Manages Profile API integration
   - Handles worker service integration
   - Implements task distribution
   - Collects worker results
   - Tracks service status
   - Monitors worker health

4. **Error Handler**

   - Detects and classifies errors
   - Implements recovery procedures
   - Manages retry policies
   - Tracks error metrics
   - Generates error alerts
   - Maintains error documentation

5. **Monitoring System**
   - Collects queue metrics
   - Tracks message metrics
   - Monitors error rates
   - Measures performance
   - Checks health status
   - Reports resource usage
   - Aggregates custom metrics

### Message Flow

1. **Profile API to Queue**

   - API sends request to queue
   - Queue validates message
   - Queue routes to appropriate worker
   - Queue tracks message status

2. **Queue to Worker**

   - Queue distributes tasks
   - Worker processes task
   - Worker sends result
   - Queue updates status

3. **Worker to Queue**

   - Worker sends result
   - Queue validates result
   - Queue updates status
   - Queue notifies API

4. **Queue to Profile API**
   - Queue sends result
   - API processes result
   - API updates status
   - API notifies client

## Technical Implementation

### Queue Configuration

```yaml
queues:
  profile_requests:
    name: profile.requests
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 3600000
      x-max-length: 10000
      x-overflow: reject-publish

  email_tasks:
    name: profile.email.tasks
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 3600000
      x-max-length: 10000
      x-overflow: reject-publish

  image_tasks:
    name: profile.image.tasks
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 3600000
      x-max-length: 10000
      x-overflow: reject-publish

  results:
    name: profile.results
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 3600000
      x-max-length: 10000
      x-overflow: reject-publish

  errors:
    name: profile.errors
    durable: true
    auto_delete: false
    arguments:
      x-message-ttl: 86400000
      x-max-length: 10000
      x-overflow: reject-publish
```

### Exchange Configuration

```yaml
exchanges:
  profile:
    name: profile
    type: topic
    durable: true
    auto_delete: false
    internal: false
    arguments: {}

  dead_letter:
    name: profile.dead-letter
    type: direct
    durable: true
    auto_delete: false
    internal: false
    arguments: {}
```

### Binding Configuration

```yaml
bindings:
  profile_requests:
    queue: profile.requests
    exchange: profile
    routing_key: profile.request.*
    arguments: {}

  email_tasks:
    queue: profile.email.tasks
    exchange: profile
    routing_key: profile.email.*
    arguments: {}

  image_tasks:
    queue: profile.image.tasks
    exchange: profile
    routing_key: profile.image.*
    arguments: {}

  results:
    queue: profile.results
    exchange: profile
    routing_key: profile.result.*
    arguments: {}

  errors:
    queue: profile.errors
    exchange: profile
    routing_key: profile.error.*
    arguments: {}
```

### Message Structure

```json
{
  "id": "string",
  "type": "string",
  "priority": "integer",
  "timestamp": "string",
  "correlation_id": "string",
  "reply_to": "string",
  "headers": {
    "key": "value"
  },
  "body": {
    "key": "value"
  }
}
```

### Error Handling

1. **Error Types**

   - Validation errors
   - Processing errors
   - Routing errors
   - Delivery errors
   - System errors

2. **Retry Policy**

   - Max retries: 3
   - Initial delay: 1s
   - Max delay: 30s
   - Backoff factor: 2

3. **Error Recovery**
   - Automatic retry
   - Manual intervention
   - Error reporting
   - Status tracking

### Monitoring

1. **Metrics**

   - Queue depth
   - Message rate
   - Error rate
   - Processing time
   - Resource usage
   - Health status

2. **Alerts**

   - Queue depth threshold
   - Error rate threshold
   - Processing time threshold
   - Resource usage threshold
   - Health check failure

3. **Logging**
   - Request logs
   - Error logs
   - Performance logs
   - Health logs
   - Resource logs
   - Custom logs

## Dependencies

- RabbitMQ Server
- Profile API Service
- Worker Email Service
- Worker Image Service
- Monitoring System
- Logging System

## Configuration

### Environment Variables

```yaml
RABBITMQ_HOST: "localhost"
RABBITMQ_PORT: "5672"
RABBITMQ_USER: "guest"
RABBITMQ_PASSWORD: "guest"
RABBITMQ_VHOST: "/"

QUEUE_PREFIX: "profile"
QUEUE_DURABLE: "true"
QUEUE_AUTO_DELETE: "false"

MESSAGE_TTL: "3600000"
MAX_QUEUE_LENGTH: "10000"
OVERFLOW_BEHAVIOR: "reject-publish"

RETRY_MAX_ATTEMPTS: "3"
RETRY_INITIAL_DELAY: "1000"
RETRY_MAX_DELAY: "30000"
RETRY_BACKOFF_FACTOR: "2"

METRICS_ENABLED: "true"
LOGGING_ENABLED: "true"
TRACING_ENABLED: "true"
```

### Kubernetes Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: profile-queue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: profile-queue
  template:
    metadata:
      labels:
        app: profile-queue
    spec:
      containers:
        - name: profile-queue
          image: profile-queue:latest
          ports:
            - containerPort: 5672
          env:
            - name: RABBITMQ_HOST
              valueFrom:
                configMapKeyRef:
                  name: profile-queue-config
                  key: RABBITMQ_HOST
            - name: RABBITMQ_PORT
              valueFrom:
                configMapKeyRef:
                  name: profile-queue-config
                  key: RABBITMQ_PORT
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "200m"
          livenessProbe:
            tcpSocket:
              port: 5672
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 5672
            initialDelaySeconds: 5
            periodSeconds: 5
```

## Development

### Prerequisites

- Go 1.21 or later
- RabbitMQ 3.12 or later
- Docker
- Kubernetes
- Make

### Building

```bash
# Build the service
make build

# Build the Docker image
make docker-build

# Push the Docker image
make docker-push
```

### Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run all tests
make test-all
```

### Deployment

```bash
# Deploy to Kubernetes
make k8s-deploy

# Undeploy from Kubernetes
make k8s-undeploy
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
