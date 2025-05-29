# Image Worker Service

-> IMPORTANT: Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.

-> LLM GUIDELINES: This document follows the LLM-friendly format. For comprehensive guidelines on LLM integration and documentation standards, refer to [LLM Guidelines](../../reference-materials/development/llm-guidelines.md).

-> CONSIDERER BEFORE UPDATING THIS FILE: The changes in this file need to be incremental, so update tasks as completed our add things that were done but do not remove any of the future plans - only change or insert new things, but do not delete. If something needs to be removed or changed, add a note instead.

-> WHERE TO GET INFORMATION TO IMPROVE THE CONTEXT: Check the `reference-materials` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions and updates to the development plan. Remember to update tasks incrementally and document all changes.

## Primary Purpose

The Image Worker Service is responsible for processing image generation tasks asynchronously. It consumes messages from the RabbitMQ queue, generates images using AI services, stores them, and updates profile information with the generated image URLs.

## Architecture

### Core Components

1. **Message Consumer**

   ```go
   // File: internal/consumer/consumer.go
   type ImageConsumer struct {
       conn    *amqp.Connection
       channel *amqp.Channel
       config  *config.Config
       ai      *ai.Generator
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

2. **AI Image Generator**

   ```go
   // File: internal/ai/generator.go
   type ImageGenerator struct {
       client  *http.Client
       config  *config.AIConfig
       metrics *metrics.ImageMetrics
   }
   ```

   Features:

   - AI service integration
   - Image generation parameters
   - Retry mechanism for failed generations
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

   - Profile image URL updates
   - Image metadata management
   - Error handling
   - Retry mechanism
   - Request tracking

4. **Metrics and Monitoring**

   ```go
   // File: internal/metrics/metrics.go
   type ImageMetrics struct {
       ProcessedMessages    int64
       FailedMessages      int64
       ImagesGenerated     int64
       GenerationFailed    int64
       ProcessingTime      time.Duration
       LastProcessedTime   time.Time
   }
   ```

   - Message processing metrics
   - Image generation metrics
   - Error tracking
   - Performance monitoring
   - Health checks

### Service Dependencies

1. **RabbitMQ**

   - Message queue for image tasks
   - Dead letter queue for failed messages
   - Message persistence
   - Queue monitoring

2. **AI Service**

   - Image generation
   - Style transfer
   - Error handling
   - Rate limiting

3. **Profile Storage Service**
   - Profile image URL updates
   - Image metadata storage
   - Error handling
   - Retry mechanism

## Implementation Details

### Project Structure

```
worker-image/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── consumer/        # Message consumer
│   ├── ai/             # AI service integration
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
    AI       AIConfig
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

type AIConfig struct {
    APIKey      string
    APIEndpoint string
    Model       string
    MaxRetries  int
    Timeout     time.Duration
}
```

Environment variables:

```bash
# RabbitMQ Configuration
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=profile-image-queue
RABBITMQ_VHOST=/

# AI Service Configuration
AI_API_KEY=your-api-key
AI_API_ENDPOINT=https://api.example.com/v1
AI_MODEL=stable-diffusion-v1
AI_MAX_RETRIES=3
AI_TIMEOUT=30s

# Storage Configuration
STORAGE_URL=http://profile-storage
STORAGE_TIMEOUT=5s
STORAGE_MAX_RETRIES=3
```

### Message Processing

1. **Message Structure**

   ```protobuf
   // File: internal/proto/image.proto
   message ImageGenerationMessage {
       string profile_id = 1;
       string prompt = 2;
       string style = 3;
       int64 created_at = 4;
   }
   ```

2. **Processing Flow**

   ```go
   // File: internal/worker/worker.go
   func (w *Worker) ProcessMessage(msg *amqp.Delivery) error {
       // 1. Parse message
       // 2. Generate image
       // 3. Store image
       // 4. Update profile
       // 5. Acknowledge message
   }
   ```

### Error Handling

1. **Message Processing Errors**

   - Invalid message format
   - Processing failures
   - Storage errors
   - AI service failures

2. **Recovery Mechanisms**
   - Message retries
   - Dead letter queue
   - Error tracking
   - Alerting

### Monitoring

1. **Health Checks**

   - Queue connection
   - AI service connection
   - Storage connection
   - Worker status

2. **Metrics**
   - Message processing rate
   - Image generation rate
   - Error rates
   - Processing time
   - Queue depth

## Development

### Prerequisites

- Go 1.21 or later
- Docker
- RabbitMQ
- AI service access
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
docker build -t worker-image .
docker run worker-image
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

2. **AI Service Security**

   - API key management
   - Rate limiting
   - Input validation
   - Error handling

3. **Storage Security**
   - HTTPS for storage service
   - Access control
   - Error handling
   - Data validation

## Best Practices

1. **Message Processing**

   - Use persistent messages
   - Implement retry mechanisms
   - Handle dead letters
   - Monitor queue health
   - Track message lifecycle

2. **Image Generation**

   - Validate prompts
   - Implement rate limiting
   - Handle generation failures
   - Track generation status
   - Monitor quality metrics

3. **Error Handling**
   - Implement retry policies
   - Use dead letter queues
   - Track error metrics
   - Log error details
   - Monitor error rates

## Known Issues and Limitations

### 1. Technical Limitations

- AI service rate limits
- Generation time constraints
- Image size limits
- Queue message size limits

### 2. Process Improvements

- Need better error tracking
- Need more comprehensive metrics
- Need improved monitoring
- Need better documentation

## Future Improvements

### 1. Short-term Goals

- Implement image validation
- Add more metrics
- Improve error handling
- Add more tests

### 2. Medium-term Goals

- Add image optimization
- Implement style transfer
- Add quality checks
- Improve monitoring

### 3. Long-term Goals

- Add image analytics
- Implement A/B testing
- Add more AI providers
- Scale infrastructure

## Version History

### Tasks History

- Initial setup
- Basic implementation
- Added metrics
- Added monitoring
