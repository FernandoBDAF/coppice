# Monitoring API Service

INITIAL CONTEXT FOR LLM - never change the context-----------------------------
-> THIS SECTION IS A GUIDELINE TO THE LLM CONSIDER BEFORE WORKING IN THIS FILE, DO NOT CHANGE THIS

-> GOES OF THE README FILE:

- This file serves as the technical documentation of the service, providing a comprehensive overview of the codebase
- It should document:
  - Service architecture and design decisions
  - Component structure and relationships
  - API endpoints and interfaces
  - Dependencies and integration points
  - Configuration and deployment details
- This is the primary reference for understanding the technical implementation
- This file should be in sync with the `/TRACKER&MANAGER.md` where development progress and tasks are tracked
- While TRACKER&MANAGER.md focuses on "what" and "when", this file focuses on "how" and "why"

-> CONSIDERER BEFORE UPDATING THIS FILE:

- Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.
- The changes in this file need to be incremental or to update informations that you confidentilly have knowlegde, they should not be guesses. If there are questions or uncertanty add comments asking for clarification instead.
- Check the `../../microservices/docs` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions, haver better context and updates to the development plan. Always compare the implementation of this project with the plan described in the docs and whenever there are inconsistancies, add comments.
- Consider structuring this documentation separating the components and describing each one of them and adding sections to how they interact with each other - because this will be very dinamic and updated during the development process it will make clear what to update after each change
- This documentation is focusing in the current component that is part of a larger project, to have a more sistemic view check `../../microservices/TRACKER&MANAGER` and `../../microservices/README`
- Do not forget to be LLM focus, so because this will be used
- For LLM-specific guidelines and patterns, refer to [LLM Integration Guide](../../docs/llm/README.md)

## Primary Purpose

The Monitoring API Service provides a centralized HTTP interface for collecting, aggregating, and exposing metrics across the microservices architecture. It handles health checks, alerts, and monitoring data collection for all services, including the Profile API, Storage API, Cache API, and Worker Services.

## Architecture

### Core Components

1. **Monitoring Manager**

   ```go
   // File: internal/monitoring/manager.go
   type MonitoringManager struct {
       prometheus *prometheus.Registry
       config     *config.Config
       logger     *logger.Logger
       metrics    *metrics.Collector
   }
   ```

   The service implements a robust monitoring manager with:

   - Metrics collection and aggregation
   - Health check management
   - Alert rule processing
   - Data retention policies
   - Error handling
   - Graceful shutdown support

2. **API Handler**

   ```go
   // File: internal/handler/api.go
   type APIHandler struct {
       monitoring *monitoring.Manager
       metrics    *metrics.Collector
       logger     *logger.Logger
   }
   ```

   Handles HTTP operations:

   - Request validation
   - Response formatting
   - Error handling
   - Metrics collection
   - Rate limiting
   - Authentication/Authorization

3. **Alert Manager**

   ```go
   // File: internal/alert/manager.go
   type AlertManager struct {
       config  *config.Config
       metrics *metrics.Collector
       logger  *logger.Logger
   }
   ```

   Manages alert operations:

   - Alert rule evaluation
   - Alert notification
   - Alert history
   - Alert status tracking
   - Alert aggregation

4. **Health Check Manager**

   ```go
   // File: internal/health/manager.go
   type HealthManager struct {
       config  *config.Config
       metrics *metrics.Collector
       logger  *logger.Logger
   }
   ```

   Manages health checks:

   - Service health monitoring
   - Dependency health checks
   - Resource usage monitoring
   - Performance metrics
   - Health status aggregation

### Service Dependencies

1. **External Services**

   - Prometheus
   - Alert Manager
   - Auth Service
   - Notification Service

2. **Internal Components**
   - Monitoring Manager
   - API Handler
   - Alert Manager
   - Health Check Manager
   - Configuration Manager

## Implementation Status

### Current State

1. **Monitoring Management**

   - [ ] Prometheus integration
   - [ ] Metrics collection
   - [ ] Health checks
   - [ ] Alert rules
   - [ ] Data retention

2. **API Endpoints**

   - [ ] Metrics endpoints
   - [ ] Health endpoints
   - [ ] Alert endpoints
   - [ ] Query endpoints
   - [ ] Management endpoints

3. **Alert System**

   - [ ] Rule evaluation
   - [ ] Notifications
   - [ ] History tracking
   - [ ] Status management
   - [ ] Aggregation

4. **Health System**
   - [ ] Service monitoring
   - [ ] Dependency checks
   - [ ] Resource monitoring
   - [ ] Status aggregation
   - [ ] Reporting

### Implementation Plan

1. **Phase 1: Core Infrastructure (Week 1)**

   - [ ] Project Structure
     - [ ] Directory setup
     - [ ] Configuration
     - [ ] Logging
     - [ ] Metrics
   - [ ] Prometheus Integration
     - [ ] Client setup
     - [ ] Metrics collection
     - [ ] Query handling
     - [ ] Error handling

2. **Phase 2: API Implementation (Week 2)**

   - [ ] HTTP Server
     - [ ] Metrics endpoints
     - [ ] Health endpoints
     - [ ] Alert endpoints
     - [ ] Error handling
   - [ ] Monitoring Operations
     - [ ] Metrics submission
     - [ ] Health updates
     - [ ] Alert management
     - [ ] Query handling

3. **Phase 3: Advanced Features (Week 3)**
   - [ ] Alert System
     - [ ] Rule engine
     - [ ] Notifications
     - [ ] History
     - [ ] Error handling
   - [ ] Health System
     - [ ] Service monitoring
     - [ ] Dependency checks
     - [ ] Resource monitoring
     - [ ] Status aggregation

## Configuration

### Environment Variables

```yaml
# Server Configuration
SERVER_PORT: 8080
SERVER_HOST: "0.0.0.0"
LOG_LEVEL: "info"

# Prometheus Configuration
PROMETHEUS_HOST: "prometheus"
PROMETHEUS_PORT: 9090
PROMETHEUS_TIMEOUT: 30s
PROMETHEUS_RETENTION: "15d"

# Alert Configuration
ALERT_RULES_PATH: "/etc/monitoring/rules"
ALERT_CHECK_INTERVAL: "1m"
ALERT_NOTIFICATION_TIMEOUT: "30s"

# Monitoring Configuration
METRICS_PORT: 9090
HEALTH_CHECK_INTERVAL: 30s
METRICS_RETENTION: "30d"
METRICS_AGGREGATION: "5m"

# API Configuration
API_RATE_LIMIT: 1000
API_TIMEOUT: 30s
API_MAX_BATCH_SIZE: 100
```

## API Endpoints

### 1. Metrics Endpoints

```http
POST /api/v1/metrics/{service}    # Submit metrics
GET /api/v1/metrics/{service}     # Query metrics
GET /api/v1/metrics/summary       # Get metrics summary
```

### 2. Health Endpoints

```http
GET /api/v1/health/{service}      # Get service health
POST /api/v1/health/{service}     # Update service health
GET /api/v1/health/summary        # Get health summary
```

### 3. Alert Endpoints

```http
GET /api/v1/alerts                # List active alerts
POST /api/v1/alerts               # Create alert rule
PUT /api/v1/alerts/{id}           # Update alert rule
DELETE /api/v1/alerts/{id}        # Delete alert rule
```

### 4. Query Endpoints

```http
GET /api/v1/query                 # Query metrics
GET /api/v1/query/range          # Query metrics range
GET /api/v1/query/series         # Query metric series
```

## Monitoring

### Metrics

1. **System Metrics**

   - CPU usage
   - Memory usage
   - Disk I/O
   - Network I/O
   - Process metrics

2. **Service Metrics**

   - Request rate
   - Response time
   - Error rate
   - Resource usage
   - Custom metrics

3. **Alert Metrics**
   - Alert rate
   - Alert duration
   - Notification rate
   - Rule evaluation time
   - Alert status

### Health Checks

1. **Service Health**

   - Service status
   - Response time
   - Error rate
   - Resource usage
   - Dependency status

2. **System Health**

   - Resource usage
   - Performance metrics
   - Error rate
   - Alert status
   - Overall status

3. **Dependency Health**
   - Database status
   - Cache status
   - Queue status
   - External service status
   - Network status

## Error Handling

### Error Types

1. **Monitoring Errors**

   - Collection errors
   - Storage errors
   - Query errors
   - Alert errors
   - Health check errors

2. **API Errors**
   - Validation errors
   - Authentication errors
   - Rate limit errors
   - System errors
   - Timeout errors

### Recovery Strategies

1. **Monitoring Recovery**

   - Collection retry
   - Storage recovery
   - Query retry
   - Alert recovery
   - Health check retry

2. **API Recovery**
   - Request retry
   - Error logging
   - Status updates
   - Alert generation
   - Circuit breaking

## Security

### Authentication

1. **API Authentication**

   - JWT tokens
   - API keys
   - Service accounts
   - Role-based access

2. **Prometheus Authentication**
   - Basic auth
   - TLS
   - Access control
   - Audit logging

### Authorization

1. **API Authorization**

   - Endpoint access
   - Resource access
   - Operation access
   - Data access

2. **Monitoring Authorization**
   - Metrics access
   - Alert access
   - Health check access
   - Query access

## Cross-References

- [Monitoring Service Patterns](../../reference-materials/development/patterns/monitoring-service-patterns.md)
- [API Patterns](../../reference-materials/development/patterns/api-patterns.md)
- [Alert Patterns](../../reference-materials/development/patterns/alert-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Prometheus Patterns](../../reference-materials/development/patterns/prometheus-patterns.md)

## Service Integration

### Integration Patterns

1. **Metrics Collection**

   ```go
   // File: internal/client/metrics.go
   type MetricsClient struct {
       httpClient *http.Client
       baseURL    string
       config     *config.Config
       logger     *logger.Logger
   }

   // Example usage in other services
   func (c *MetricsClient) SubmitMetrics(metrics *Metrics) error {
       // Submit service metrics
       // - Request rates
       // - Response times
       // - Error rates
       // - Resource usage
       // - Custom metrics
   }
   ```

2. **Health Reporting**

   ```go
   // File: internal/client/health.go
   type HealthClient struct {
       httpClient *http.Client
       baseURL    string
       config     *config.Config
       logger     *logger.Logger
   }

   // Example usage in other services
   func (c *HealthClient) UpdateHealth(status *HealthStatus) error {
       // Update service health status
       // - Service status
       // - Dependency status
       // - Resource usage
       // - Performance metrics
   }
   ```

3. **Alert Integration**

   ```go
   // File: internal/client/alert.go
   type AlertClient struct {
       httpClient *http.Client
       baseURL    string
       config     *config.Config
       logger     *logger.Logger
   }

   // Example usage in other services
   func (c *AlertClient) CreateAlert(alert *Alert) error {
       // Create service alerts
       // - Error alerts
       // - Performance alerts
       // - Resource alerts
       // - Custom alerts
   }
   ```

### Integration Examples

1. **Profile API Integration**

   ```go
   // File: services/profile-api/internal/monitoring/client.go
   type MonitoringClient struct {
       metrics *MetricsClient
       health  *HealthClient
       alert   *AlertClient
   }

   func (c *MonitoringClient) Initialize() error {
       // Initialize monitoring clients
       // - Configure metrics collection
       // - Set up health reporting
       // - Configure alerts
   }

   func (c *MonitoringClient) TrackRequest(req *http.Request) {
       // Track API request
       // - Record request metrics
       // - Update health status
       // - Check for alerts
   }
   ```

2. **Storage API Integration**

   ```go
   // File: services/storage-api/internal/monitoring/client.go
   type MonitoringClient struct {
       metrics *MetricsClient
       health  *HealthClient
       alert   *AlertClient
   }

   func (c *MonitoringClient) TrackOperation(op *StorageOperation) {
       // Track storage operation
       // - Record operation metrics
       // - Update health status
       // - Check for alerts
   }
   ```

3. **Cache API Integration**

   ```go
   // File: services/cache-api/internal/monitoring/client.go
   type MonitoringClient struct {
       metrics *MetricsClient
       health  *HealthClient
       alert   *AlertClient
   }

   func (c *MonitoringClient) TrackCacheOperation(op *CacheOperation) {
       // Track cache operation
       // - Record cache metrics
       // - Update health status
       // - Check for alerts
   }
   ```

### Integration Guidelines

1. **Metrics Collection**

   - Use consistent metric names
   - Include service labels
   - Set appropriate aggregation
   - Configure retention policies
   - Handle collection errors

2. **Health Reporting**

   - Report service status
   - Include dependency status
   - Track resource usage
   - Monitor performance
   - Handle reporting errors

3. **Alert Management**

   - Define alert rules
   - Configure notifications
   - Track alert history
   - Handle alert errors
   - Manage alert lifecycle

4. **Best Practices**

   - Use connection pooling
   - Implement retry logic
   - Handle timeouts
   - Manage resources
   - Follow security guidelines

### Configuration

1. **Client Configuration**

   ```yaml
   # Monitoring Client Configuration
   MONITORING_API_HOST: "monitoring-api"
   MONITORING_API_PORT: 8080
   MONITORING_API_TIMEOUT: 30s
   MONITORING_API_RETRY_ATTEMPTS: 3
   MONITORING_API_RETRY_DELAY: 1s
   ```

2. **Metrics Configuration**

   ```yaml
   # Metrics Configuration
   METRICS_ENABLED: true
   METRICS_INTERVAL: 15s
   METRICS_BATCH_SIZE: 100
   METRICS_RETENTION: 30d
   METRICS_AGGREGATION: 5m
   ```

3. **Health Configuration**

   ```yaml
   # Health Configuration
   HEALTH_CHECK_ENABLED: true
   HEALTH_CHECK_INTERVAL: 30s
   HEALTH_CHECK_TIMEOUT: 5s
   HEALTH_CHECK_RETRY_ATTEMPTS: 3
   ```

4. **Alert Configuration**

   ```yaml
   # Alert Configuration
   ALERT_ENABLED: true
   ALERT_CHECK_INTERVAL: 1m
   ALERT_NOTIFICATION_TIMEOUT: 30s
   ALERT_RETRY_ATTEMPTS: 3
   ```

### Error Handling

1. **Client Errors**

   ```go
   // File: internal/client/errors.go
   type ClientError struct {
       Code    string
       Message string
       Cause   error
   }

   func (e *ClientError) Error() string {
       return fmt.Sprintf("%s: %s", e.Code, e.Message)
   }
   ```

2. **Recovery Strategies**

   - Implement retry logic
   - Use circuit breakers
   - Handle timeouts
   - Log errors
   - Report failures

3. **Error Types**

   - Connection errors
   - Timeout errors
   - Validation errors
   - Authentication errors
   - Rate limit errors

### Security

1. **Authentication**

   ```go
   // File: internal/client/auth.go
   type AuthConfig struct {
       APIKey     string
       ServiceID  string
       Role       string
       Permissions []string
   }
   ```

2. **Authorization**

   - Role-based access
   - Resource permissions
   - Operation permissions
   - Data access control

3. **Security Best Practices**

   - Use TLS
   - Rotate credentials
   - Audit logging
   - Access control
   - Security monitoring
