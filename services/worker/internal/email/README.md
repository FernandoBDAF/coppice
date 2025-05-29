# Email Worker Service

-> IMPORTANT: Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.

-> LLM GUIDELINES: This document follows the LLM-friendly format. For comprehensive guidelines on LLM integration and documentation standards, refer to [LLM Guidelines](../../reference-materials/development/llm-guidelines.md).

-> CONSIDERER BEFORE UPDATING THIS FILE: The changes in this file need to be incremental, so update tasks as completed our add things that were done but do not remove any of the future plans - only change or insert new things, but do not delete. If something needs to be removed or changed, add a note instead.

-> WHERE TO GET INFORMATION TO IMPROVE THE CONTEXT: Check the `reference-materials` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions and updates to the development plan. Remember to update tasks incrementally and document all changes.

## Primary Purpose

The Email Worker Service is responsible for processing email validation tasks asynchronously. It consumes messages from the RabbitMQ queue, sends validation emails, and updates profile status based on email validation results.

## Architecture

### Core Components

1. **Message Consumer**

   ```go
   // File: internal/consumer/consumer.go
   type EmailConsumer struct {
       conn    *amqp.Connection
       channel *amqp.Channel
       config  *config.Config
       email   *email.Sender
       storage *storage.Client
   }
   ```

   The service implements a robust message consumer with:

   - Connection management and recovery
   - Message acknowledgment handling
   - Error handling and retries
   - Dead letter queue integration
   - Message validation
   - Graceful shutdown support

2. **Email Sender**

   ```go
   // File: internal/email/sender.go
   type EmailSender struct {
       client  *smtp.Client
       config  *config.EmailConfig
       metrics *metrics.EmailMetrics
   }
   ```

   Features:

   - SMTP client with TLS support
   - Email template management
   - Retry mechanism for failed sends
   - Rate limiting
   - Error tracking
   - Metrics collection

3. **Storage Integration**

   ```go
   // File: internal/storage/client.go
   type StorageClient struct {
       client  *http.Client
       baseURL string
       config  *config.StorageConfig
   }
   ```

   - Profile status updates
   - Validation token management
   - Error handling
   - Retry mechanism
   - Request tracking

4. **Metrics and Monitoring**

   ```go
   // File: internal/metrics/metrics.go
   type EmailMetrics struct {
       ProcessedMessages    int64
       FailedMessages      int64
       EmailSentCount      int64
       EmailFailedCount    int64
       ProcessingTime      time.Duration
       LastProcessedTime   time.Time
   }
   ```

   - Message processing metrics
   - Email sending metrics
   - Error tracking
   - Performance monitoring
   - Health checks

### Service Dependencies

1. **RabbitMQ**

   - Message queue for email tasks
   - Dead letter queue for failed messages
   - Message persistence
   - Queue monitoring

2. **SMTP Server**

   - Email delivery
   - TLS support
   - Rate limiting
   - Error handling

3. **Profile Storage Service**
   - Profile status updates
   - Validation token storage
   - Error handling
   - Retry mechanism

## Implementation Details

### Project Structure

```
worker-email/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── consumer/        # Message consumer
│   ├── email/          # Email sending
│   ├── metrics/        # Metrics collection
│   ├── storage/        # Storage service client
│   └── worker/         # Worker implementation
└── k8s/               # Kubernetes configuration
```

### Configuration

```go
// File: internal/config/config.go
type Config struct {
    RabbitMQ RabbitMQConfig
    Email    EmailConfig
    Storage  StorageConfig
    Metrics  MetricsConfig
}

type RabbitMQConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Queue    string
    VHost    string
}

type EmailConfig struct {
    SMTPHost     string
    SMTPPort     int
    Username     string
    Password     string
    FromAddress  string
    TemplatePath string
}
```

Environment variables:

```bash
# RabbitMQ Configuration
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=profile-email-queue
RABBITMQ_VHOST=/

# Email Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user
SMTP_PASSWORD=pass
EMAIL_FROM=noreply@example.com
EMAIL_TEMPLATE_PATH=/templates

# Storage Configuration
STORAGE_URL=http://profile-storage
STORAGE_TIMEOUT=5s
STORAGE_MAX_RETRIES=3
```

### Message Processing

1. **Message Structure**

   ```protobuf
   // File: internal/proto/email.proto
   message EmailValidationMessage {
       string profile_id = 1;
       string email = 2;
       string validation_token = 3;
       int64 created_at = 4;
   }
   ```

2. **Processing Flow**

   ```go
   // File: internal/worker/worker.go
   func (w *Worker) ProcessMessage(msg *amqp.Delivery) error {
       // 1. Parse message
       // 2. Send validation email
       // 3. Update profile status
       // 4. Acknowledge message
   }
   ```

### Error Handling

1. **Message Processing Errors**

   - Invalid message format
   - Processing failures
   - Storage errors
   - Email sending failures

2. **Recovery Mechanisms**
   - Message retries
   - Dead letter queue
   - Error tracking
   - Alerting

### Monitoring

1. **Health Checks**

   - Queue connection
   - SMTP connection
   - Storage connection
   - Worker status

2. **Metrics**
   - Message processing rate
   - Email sending rate
   - Error rates
   - Processing time
   - Queue depth

## Development

### Prerequisites

- Go 1.21 or later
- Docker
- RabbitMQ
- SMTP server
- Profile Storage service

### Setup

```bash
# Clone repository
git clone [repository-url]

# Install dependencies
go mod download

# Run locally
go run cmd/main.go

# Run with Docker
docker build -t worker-email .
docker run worker-email
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Security

1. **Message Security**

   - TLS for RabbitMQ
   - Message encryption
   - Access control
   - Audit logging

2. **Email Security**

   - TLS for SMTP
   - Credential management
   - Rate limiting
   - Spam prevention

3. **Storage Security**
   - HTTPS for storage service
   - Token validation
   - Access control
   - Error handling

## Best Practices

1. **Message Processing**

   - Use persistent messages
   - Implement retry mechanisms
   - Handle dead letters
   - Monitor queue health
   - Track message lifecycle

2. **Email Sending**

   - Use templates
   - Implement rate limiting
   - Handle bounces
   - Track delivery status
   - Monitor spam scores

3. **Error Handling**
   - Implement retry policies
   - Use dead letter queues
   - Track error metrics
   - Log error details
   - Monitor error rates

## Known Issues and Limitations

### 1. Technical Limitations

- Email delivery not guaranteed
- SMTP server dependencies
- Rate limiting constraints
- Queue message size limits

### 2. Process Improvements

- Need better error tracking
- Need more comprehensive metrics
- Need improved monitoring
- Need better documentation

## Future Improvements

### 1. Short-term Goals

- Implement email templates
- Add more metrics
- Improve error handling
- Add more tests

### 2. Medium-term Goals

- Add email tracking
- Implement bounce handling
- Add spam prevention
- Improve monitoring

### 3. Long-term Goals

- Add email analytics
- Implement A/B testing
- Add more email providers
- Scale infrastructure

## Version History

### Tasks History

- Initial setup
- Basic implementation
- Added metrics
- Added monitoring
