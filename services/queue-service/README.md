# Queue Service

## Overview

The Queue Service is a critical component of our microservices architecture that provides reliable message queuing and event processing capabilities. It handles asynchronous communication between services, ensuring message delivery, persistence, and processing.

## Role in the System

The Queue Service interacts with several components:

1. **Other Microservices**

   - Auth Service: Authentication events
   - Profile Service: Profile update events
   - Cache Service: Cache invalidation events
   - Storage Service: Data operation events
   - Worker Service: Job processing events
   - Monitoring Service: Queue metrics

2. **External Services**
   - RabbitMQ: Message broker
   - Redis: Message persistence
   - Monitoring Service: Metrics collection

## Main Functionalities

1. **Message Queue Management**

   - Queue creation and configuration
   - Message routing
   - Dead letter handling
   - Queue monitoring

2. **Message Processing**

   - Message validation
   - Message transformation
   - Message persistence
   - Message delivery

3. **Event Handling**

   - Event publishing
   - Event subscription
   - Event routing
   - Event persistence

4. **Monitoring**
   - Queue metrics
   - Message rates
   - Error tracking
   - Performance monitoring

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make
- RabbitMQ 3.9+
- Redis 7.0+

### Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd queue-service
```

2. Install dependencies:

```bash
make deps
```

3. Configure environment:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start the service:

```bash
make run
```

### Configuration

Essential environment variables:

```bash
# Service Configuration
SERVICE_NAME=queue-service
SERVICE_PORT=8080
METRICS_PORT=9090

# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Queue Configuration
QUEUE_PREFIX=profile
MAX_RETRIES=3
RETRY_DELAY=5000
MESSAGE_TTL=86400000
```

### Running with Docker

```bash
# Build and run
docker-compose up --build
```

## Development

### Common Tasks

1. Run tests:

```bash
make test
```

2. Build service:

```bash
make build
```

3. Run linter:

```bash
make lint
```

### Project Structure

```
queue-service/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── api/         # API handlers and routes
│   ├── queue/       # Queue management
│   ├── message/     # Message handling
│   ├── event/       # Event management
│   ├── config/      # Configuration
│   └── service/     # Business logic
├── pkg/             # Public libraries
├── test/            # Test files
└── docs/            # Documentation
```

## Documentation

- [Context](./CONTEXT.md) - Technical details and architecture
- [Interface](./INTERFACE.md) - Service connections and APIs
- [Tracker](./TRACKER.md) - Development tasks and progress

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
