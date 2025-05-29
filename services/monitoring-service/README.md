# Monitoring Service

## Overview

The Monitoring Service is a central component of our microservices architecture that provides comprehensive monitoring, alerting, and observability capabilities. It collects metrics, logs, and traces from all services, providing real-time insights into system health and performance.

## Role in the System

The Monitoring Service interacts with several components:

1. **Other Microservices**

   - Auth Service: Authentication metrics and health checks
   - Profile Service: User operations metrics
   - Cache Service: Cache performance metrics
   - Storage Service: Storage operations metrics
   - Worker Service: Job processing metrics
   - Queue Service: Message queue metrics

2. **External Services**
   - Prometheus: Metrics collection and storage
   - Grafana: Metrics visualization
   - ELK Stack: Log aggregation and analysis
   - Jaeger: Distributed tracing

## Main Functionalities

1. **Metrics Collection**

   - Service health metrics
   - Performance metrics
   - Business metrics
   - Resource utilization

2. **Alerting**

   - Threshold-based alerts
   - Anomaly detection
   - Alert routing
   - Alert management

3. **Logging**

   - Centralized log collection
   - Log aggregation
   - Log analysis
   - Log retention

4. **Tracing**
   - Distributed tracing
   - Request tracking
   - Performance analysis
   - Dependency mapping

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make
- Prometheus
- Grafana
- ELK Stack
- Jaeger

### Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd monitoring-service
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
SERVICE_NAME=monitoring-service
SERVICE_PORT=8080
METRICS_PORT=9090

# Prometheus Configuration
PROMETHEUS_HOST=localhost
PROMETHEUS_PORT=9090
PROMETHEUS_RETENTION=15d

# Grafana Configuration
GRAFANA_HOST=localhost
GRAFANA_PORT=3000
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=admin

# ELK Configuration
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
KIBANA_HOST=localhost
KIBANA_PORT=5601

# Jaeger Configuration
JAEGER_HOST=localhost
JAEGER_PORT=6831

# Alerting Configuration
ALERT_MANAGER_HOST=localhost
ALERT_MANAGER_PORT=9093
SLACK_WEBHOOK_URL=your-webhook-url
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
monitoring-service/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── api/         # API handlers and routes
│   ├── metrics/     # Metrics collection
│   ├── alerts/      # Alert management
│   ├── logging/     # Log management
│   ├── tracing/     # Trace management
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
