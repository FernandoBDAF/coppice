# Cache Service

## Overview

The Cache Service is a critical component of our microservices architecture that provides distributed caching capabilities to other services. It implements a Redis-based caching layer with advanced features like distributed locking, cache invalidation, and automatic expiration.

## Role in the System

The Cache Service interacts with several components:

1. **Other Microservices**

   - Auth Service: Session storage and token caching
   - Profile Service: User profile caching
   - Storage Service: File metadata caching
   - Worker Service: Job result caching

2. **External Services**
   - Redis: Primary caching backend
   - Monitoring Service: Cache metrics and health status
   - Logging Service: Cache operation logs

## Main Functionalities

1. **Distributed Caching**

   - Key-value storage
   - TTL-based expiration
   - Cache invalidation
   - Cache warming

2. **Advanced Features**

   - Distributed locking
   - Atomic operations
   - Batch operations
   - Pub/Sub messaging

3. **Monitoring**
   - Cache hit/miss rates
   - Memory usage
   - Operation latency
   - Error rates

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make
- Redis 7.0 or higher

### Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd cache-service
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
SERVICE_NAME=cache-service
SERVICE_PORT=8080
METRICS_PORT=9090

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your-password
REDIS_DB=0

# Cache Configuration
CACHE_TTL=3600
CACHE_MAX_SIZE=1000
CACHE_CLEANUP_INTERVAL=300

# Monitoring
ENABLE_METRICS=true
ENABLE_TRACING=true
LOG_LEVEL=info
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
cache-service/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── api/         # API handlers and routes
│   ├── cache/       # Cache implementation
│   ├── config/      # Configuration
│   ├── models/      # Data models
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
