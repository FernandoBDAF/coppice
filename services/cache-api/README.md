# Profile Cache Service

## Overview

The Profile Cache Service provides centralized caching capabilities for the microservices architecture. It integrates with shared libraries to provide robust caching with proper monitoring, logging, and error handling.

## Architecture

### Core Components

1. **Cache Layer**

   - Cache operations
   - Data persistence
   - Eviction management
   - Connection pooling

2. **Service Layer**

   - Business logic
   - Data transformation
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
   logger := logging.NewLogger("profile-cache")

   // Usage example
   logger.Info("Processing cache operation",
       logging.WithField("operation", "get"),
       logging.WithField("key", cacheKey))
   ```

2. **Monitoring Library**

   ```go
   // Initialize collector
   monitor := monitoring.NewCollector("profile-cache")

   // Usage example
   monitor.IncCacheOperations("get")
   defer monitor.ObserveDuration("cache_get")
   ```

3. **Cache Library**

   ```go
   // Initialize cache client
   cacheClient := cache.NewAPIClient(cache.Config{
       Endpoint: "http://cache-api:8080",
       Timeout:  time.Second * 5,
   })

   // Usage example
   value, err = cacheClient.Get(ctx, cacheKey)
   ```

### Service Integration Library

```go
// Initialize service integration
integration := integration.NewServiceIntegration(integration.Config{
    ServiceName: "profile-cache",
    Discovery:   "kubernetes",
})

// Register health checks
integration.RegisterHealthCheck("redis", func() error {
    return redis.HealthCheck(ctx)
})

// Use circuit breaker
breaker := integration.NewCircuitBreaker("redis", integration.CircuitBreakerConfig{
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

### 1. Cache Invalidation

```go
func (s *CacheService) InvalidateProfile(ctx context.Context, profileID string) error {
    // Create transaction context
    txCtx := integration.NewTransactionContext(ctx)
    defer txCtx.Cleanup()

    // Begin transaction
    if err := txCtx.Begin(); err != nil {
        return err
    }

    // Invalidate profile cache
    if err := s.cache.Delete(txCtx, "profile:"+profileID); err != nil {
        txCtx.Rollback()
        return err
    }

    // Invalidate related caches
    if err := s.cache.DeletePattern(txCtx, "profile:"+profileID+":*"); err != nil {
        txCtx.Rollback()
        return err
    }

    return txCtx.Commit()
}
```

### 2. Cache Warming

```go
func (s *CacheService) WarmCache(ctx context.Context, profiles []*Profile) error {
    // Setup batch processor
    processor := integration.NewBatchProcessor(integration.BatchProcessorConfig{
        BatchSize: 100,
        Workers:   5,
    })

    // Process profiles in batches
    return processor.Process(ctx, profiles, func(ctx context.Context, batch []*Profile) error {
        return s.cache.BatchSet(ctx, batch)
    })
}
```

## Configuration

### Base Configuration

```yaml
service:
  name: profile-cache
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
cache:
  type: redis
  host: redis
  port: 6379
  password: ${REDIS_PASSWORD}
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
    case errors.Is(err, cache.ErrNotFound):
        logger.Warn("Cache miss", logging.WithError(err))
        monitor.IncCacheMisses()
    case errors.Is(err, cache.ErrConnection):
        logger.Error("Cache connection error", logging.WithError(err))
        monitor.IncConnectionErrors()
    case errors.Is(err, cache.ErrTimeout):
        logger.Error("Cache timeout", logging.WithError(err))
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
integration.RegisterHealthCheck("redis", func() error {
    return redis.HealthCheck(ctx)
})
```

## Metrics

### Standard Metrics

```go
// Cache metrics
monitor.IncCacheOperations("get")
monitor.ObserveCacheLatency("get", latency)

// Redis metrics
monitor.IncRedisOperations("get")
monitor.ObserveRedisLatency("get", latency)
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
   go build -o profile-cache ./cmd/profile-cache
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
   k6 run ./tests/load/profile-cache.js
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
   kubectl logs -n microservices -l app=profile-cache
   ```

### Docker

1. Build image:

   ```bash
   docker build -t profile-cache:latest .
   ```

2. Run container:
   ```bash
   docker run -p 8080:8080 profile-cache:latest
   ```

## Monitoring

### Prometheus Metrics

- Cache operation rates
- Error rates
- Latency percentiles
- Connection pool metrics
- Memory usage metrics

### Grafana Dashboards

- Service overview
- Error rates
- Latency trends
- Cache performance
- Memory usage

## Logging

### Log Levels

- ERROR: Service errors
- WARN: Cache misses
- INFO: Request processing
- DEBUG: Detailed operations

### Log Fields

- service: profile-cache
- trace_id: Request tracing
- cache_key: Cache key
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

- Redis Cache
- Cache API Service

### Shared Libraries

- Logging Library
- Monitoring Library
- Cache Client Library
- Service Integration Library

## API Documentation

### OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: Profile Cache API
  version: 1.0.0
paths:
  /cache:
    get:
      summary: Get cached value
      responses:
        "200":
          description: Success
    post:
      summary: Set cached value
      responses:
        "201":
          description: Created
  /cache/{key}:
    get:
      summary: Get cached value by key
      responses:
        "200":
          description: Success
    put:
      summary: Update cached value
      responses:
        "200":
          description: Success
    delete:
      summary: Delete cached value
      responses:
        "204":
          description: No Content
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

1. **Cache Layer**

   - [ ] Cache client implementation
   - [ ] Cache configuration
   - [ ] Cache type definitions
   - [ ] Error handling
   - [ ] Metrics collection

2. **Service Layer**

   - [ ] Business logic implementation
   - [ ] Cache transformation
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

2. **Phase 2: Cache Implementation**

   - [ ] Redis client setup
   - [ ] Cache operations
   - [ ] Error handling
   - [ ] Metrics collection

3. **Phase 3: Service Integration**
   - [ ] Shared libraries integration
   - [ ] API services communication
   - [ ] Circuit breaking
   - [ ] Retry mechanisms

## API Endpoints

### 1. Cache Operations

```http
GET /api/v1/cache/{key}
POST /api/v1/cache/{key}
DELETE /api/v1/cache/{key}
```

### 2. Cache Management

```http
GET /api/v1/cache/stats
POST /api/v1/cache/flush
GET /api/v1/cache/status
```

### 3. Health and Metrics

```http
GET /health
GET /ready
GET /metrics
```

## Error Types

### 1. Cache Errors

- Connection errors
- Operation errors
- Memory errors
- Timeout errors
- Serialization errors

### 2. Service Errors

- Validation errors
- Dependency errors
- Resource errors
- Timeout errors

### Recovery Strategies

### 1. Cache Recovery

- Connection retry
- Operation retry
- Memory management
- Error logging
- Circuit breaking

### 2. Service Recovery

- Validation recovery
- Dependency recovery
- Resource recovery
- Error logging
- Circuit breaking

## Cross-References

- [Cache Service Patterns](../../reference-materials/development/patterns/cache-service-patterns.md)
- [Service Integration Patterns](../../reference-materials/development/patterns/service-integration-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Error Handling Patterns](../../reference-materials/development/patterns/error-handling-patterns.md)
