# Monitoring Base Library

INITIAL CONTEXT FOR LLM - never change the context-----------------------------
-> THIS SECTION IS A GUIDELINE TO THE LLM CONSIDER BEFORE WORKING IN THIS FILE, DO NOT CHANGE THIS

-> GOES OF THE README FILE:

- This file serves as the technical documentation of the library, providing a comprehensive overview of the codebase
- It should document:
  - Library architecture and design decisions
  - Component structure and relationships
  - Integration points and interfaces
  - Configuration and usage details
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

The Monitoring Base Library provides a standardized foundation for implementing Prometheus-based monitoring across all microservices. It offers common metrics collection, health checks, and alerting capabilities while allowing services to extend with their specific monitoring needs.

## Architecture

### Core Components

1. **Base Collector**

   ```go
   // File: collector.go
   type BaseCollector struct {
       serviceName string
       registry    *prometheus.Registry
       metrics     map[string]prometheus.Collector
       logger      *logger.Logger
       config      *config.Config
   }

   func NewBaseCollector(serviceName string) *BaseCollector {
       return &BaseCollector{
           serviceName: serviceName,
           registry:    prometheus.NewRegistry(),
           metrics:     make(map[string]prometheus.Collector),
           logger:      logger.New(),
           config:      config.New(),
       }
   }
   ```

2. **Standard Metrics**

   ```go
   // File: metrics.go
   var (
       // HTTP Metrics
       httpRequestsTotal = prometheus.NewCounterVec(
           prometheus.CounterOpts{
               Name: "http_requests_total",
               Help: "Total number of HTTP requests",
           },
           []string{"service", "endpoint", "method", "status"},
       )

       httpRequestDuration = prometheus.NewHistogramVec(
           prometheus.HistogramOpts{
               Name: "http_request_duration_seconds",
               Help: "HTTP request duration in seconds",
           },
           []string{"service", "endpoint", "method"},
       )

       // System Metrics
       systemMemoryUsage = prometheus.NewGaugeVec(
           prometheus.GaugeOpts{
               Name: "system_memory_usage_bytes",
               Help: "System memory usage in bytes",
           },
           []string{"service"},
       )

       // Error Metrics
       errorTotal = prometheus.NewCounterVec(
           prometheus.CounterOpts{
               Name: "error_total",
               Help: "Total number of errors",
           },
           []string{"service", "type", "code"},
       )
   )
   ```

3. **Health Checker**

   ```go
   // File: health.go
   type HealthChecker struct {
       serviceName string
       checks      map[string]Check
       status      *atomic.Value
       logger      *logger.Logger
   }

   type Check struct {
       Name     string
       Check    func() error
       Interval time.Duration
   }
   ```

4. **Alert Manager**

   ```go
   // File: alert.go
   type AlertManager struct {
       serviceName string
       rules       map[string]AlertRule
       logger      *logger.Logger
       config      *config.Config
   }

   type AlertRule struct {
       Name        string
       Expr        string
       Duration    string
       Labels      map[string]string
       Annotations map[string]string
   }
   ```

## Service Integration

### Basic Integration

```go
// File: services/profile-api/internal/monitoring/collector.go
package monitoring

import (
    "github.com/your-org/monitoring"
)

type ProfileCollector struct {
    *monitoring.BaseCollector
    customMetrics map[string]prometheus.Collector
}

func NewProfileCollector() *ProfileCollector {
    base := monitoring.NewBaseCollector("profile-api")
    return &ProfileCollector{
        BaseCollector: base,
        customMetrics: make(map[string]prometheus.Collector),
    }
}

func (c *ProfileCollector) Initialize() error {
    // Register standard metrics
    c.RegisterStandardMetrics()

    // Register custom metrics
    c.RegisterCustomMetrics()

    // Initialize health checks
    c.InitializeHealthChecks()

    return nil
}
```

### Custom Metrics Implementation

```go
// File: services/profile-api/internal/monitoring/metrics.go
package monitoring

var (
    // Profile-specific metrics
    profileCreationTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "profile_creation_total",
            Help: "Total number of profile creations",
        },
        []string{"service", "status"},
    )

    profileValidationDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "profile_validation_duration_seconds",
            Help: "Profile validation duration in seconds",
        },
        []string{"service", "type"},
    )
)

func (c *ProfileCollector) RegisterCustomMetrics() {
    c.customMetrics["profile_creation_total"] = profileCreationTotal
    c.customMetrics["profile_validation_duration"] = profileValidationDuration

    for _, metric := range c.customMetrics {
        c.registry.MustRegister(metric)
    }
}
```

### Health Check Implementation

```go
// File: services/profile-api/internal/monitoring/health.go
package monitoring

func (c *ProfileCollector) InitializeHealthChecks() {
    // Add standard health checks
    c.AddHealthCheck("database", c.checkDatabase, 30*time.Second)
    c.AddHealthCheck("cache", c.checkCache, 30*time.Second)

    // Add custom health checks
    c.AddHealthCheck("profile-storage", c.checkProfileStorage, 30*time.Second)
    c.AddHealthCheck("email-service", c.checkEmailService, 30*time.Second)
}

func (c *ProfileCollector) checkDatabase() error {
    // Implement database health check
    return nil
}

func (c *ProfileCollector) checkCache() error {
    // Implement cache health check
    return nil
}
```

### Alert Rule Implementation

```go
// File: services/profile-api/internal/monitoring/alerts.go
package monitoring

func (c *ProfileCollector) InitializeAlertRules() {
    // Add standard alert rules
    c.AddAlertRule(monitoring.AlertRule{
        Name:     "high_error_rate",
        Expr:     "rate(error_total{service=\"profile-api\"}[5m]) > 0.1",
        Duration: "5m",
        Labels: map[string]string{
            "severity": "critical",
        },
        Annotations: map[string]string{
            "summary": "High error rate detected",
        },
    })

    // Add custom alert rules
    c.AddAlertRule(monitoring.AlertRule{
        Name:     "profile_creation_failure",
        Expr:     "rate(profile_creation_total{status=\"error\"}[5m]) > 0",
        Duration: "5m",
        Labels: map[string]string{
            "severity": "warning",
        },
        Annotations: map[string]string{
            "summary": "Profile creation failures detected",
        },
    })
}
```

## Configuration

### Base Configuration

```yaml
# File: config/base.yaml
monitoring:
  base:
    enabled: true
    scrape_interval: 15s
    metrics_path: /metrics
    timeout: 10s

  metrics:
    http:
      enabled: true
      buckets: [0.1, 0.5, 1, 2, 5]
    system:
      enabled: true
      interval: 60s
    error:
      enabled: true

  health:
    enabled: true
    interval: 30s
    timeout: 5s

  alert:
    enabled: true
    check_interval: 1m
    notification_timeout: 30s
```

### Service-Specific Configuration

```yaml
# File: services/profile-api/config/monitoring.yaml
monitoring:
  custom:
    enabled: true

  metrics:
    profile:
      enabled: true
      buckets: [0.1, 0.5, 1, 2, 5]
    validation:
      enabled: true
      buckets: [0.1, 0.5, 1, 2, 5]
    cache:
      enabled: true

  health:
    database:
      enabled: true
      interval: 30s
    cache:
      enabled: true
      interval: 30s
    profile-storage:
      enabled: true
      interval: 30s

  alert:
    profile_creation:
      enabled: true
      threshold: 0
      duration: 5m
    validation_failure:
      enabled: true
      threshold: 0.1
      duration: 5m
```

## Usage Examples

### HTTP Middleware

```go
// File: middleware/metrics.go
package middleware

import (
    "github.com/your-org/monitoring"
)

func MetricsMiddleware(collector *monitoring.BaseCollector) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Track request
            httpRequestsTotal.WithLabelValues(
                collector.ServiceName(),
                r.URL.Path,
                r.Method,
            ).Inc()

            next.ServeHTTP(w, r)

            // Track duration
            httpRequestDuration.WithLabelValues(
                collector.ServiceName(),
                r.URL.Path,
                r.Method,
            ).Observe(time.Since(start).Seconds())
        })
    }
}
```

### Custom Metric Collection

```go
// File: services/profile-api/internal/handler/profile.go
package handler

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // Create profile
    profile, err := h.service.CreateProfile(r.Context(), req)
    if err != nil {
        // Track error
        errorTotal.WithLabelValues(
            h.collector.ServiceName(),
            "profile_creation",
            "500",
        ).Inc()

        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Track success
    profileCreationTotal.WithLabelValues(
        h.collector.ServiceName(),
        "success",
    ).Inc()

    // Track duration
    profileValidationDuration.WithLabelValues(
        h.collector.ServiceName(),
        "creation",
    ).Observe(time.Since(start).Seconds())
}
```

## Best Practices

1. **Metric Naming**

   - Use consistent naming conventions
   - Include service name in metric names
   - Use appropriate metric types
   - Add helpful descriptions
   - Use appropriate labels

2. **Health Checks**

   - Check all dependencies
   - Set appropriate intervals
   - Handle timeouts
   - Provide detailed status
   - Include custom checks

3. **Alert Rules**

   - Set appropriate thresholds
   - Use meaningful durations
   - Include helpful annotations
   - Set proper severity levels
   - Test alert conditions

4. **Performance**
   - Use appropriate buckets
   - Set reasonable intervals
   - Handle errors gracefully
   - Monitor resource usage
   - Use efficient collection

## Cross-References

- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Prometheus Patterns](../../reference-materials/development/patterns/prometheus-patterns.md)
- [Alert Patterns](../../reference-materials/development/patterns/alert-patterns.md)
- [Health Check Patterns](../../reference-materials/development/patterns/health-check-patterns.md)
- [Metric Patterns](../../reference-materials/development/patterns/metric-patterns.md)
