# Worker Service

## Overview

The Worker Service is a critical component of our microservices architecture that handles background job processing, scheduled tasks, and asynchronous operations. It provides reliable job execution, monitoring, and management capabilities for various tasks across the system.

## Role in the System

The Worker Service interacts with several components:

1. **Other Microservices**

   - Auth Service: Authentication and authorization
   - Profile Service: Profile-related tasks
   - Cache Service: Cache invalidation tasks
   - Storage Service: Data processing tasks
   - Queue Service: Job queue management
   - Monitoring Service: Job metrics and health

2. **External Services**
   - RabbitMQ: Job queue
   - Redis: Job state and locks
   - Monitoring Service: Metrics collection

## Main Functionalities

1. **Job Processing**

   - Background job execution
   - Task scheduling
   - Job prioritization
   - Progress tracking
   - Error handling

2. **Task Management**

   - Task creation and configuration
   - Task scheduling
   - Task monitoring
   - Task cancellation
   - Task retry

3. **Worker Types**

   - Email validation worker
   - Image generation worker
   - Data processing worker
   - Cache invalidation worker
   - Profile update worker

4. **Monitoring**
   - Job metrics
   - Worker health
   - Performance tracking
   - Error monitoring
   - Resource utilization

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
cd worker-service
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
SERVICE_NAME=worker-service
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

# Worker Configuration
MAX_WORKERS=10
JOB_TIMEOUT=3600
RETRY_ATTEMPTS=3
RETRY_DELAY=5000
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
worker-service/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── api/         # API handlers and routes
│   ├── worker/      # Worker implementations
│   ├── job/         # Job management
│   ├── scheduler/   # Task scheduling
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
