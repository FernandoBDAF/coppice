# Shared Components and Best Practices

## Overview

This document outlines the shared components and best practices that should be implemented across all microservices to ensure consistency, reduce code duplication, and improve maintainability.

## Shared Components

### 1. Logging Package

```go
// pkg/logging/logger.go
package logging

import (
    "context"
    "log"
    "os"
)

type Logger struct {
    *log.Logger
    level string
}

func NewLogger(level string) *Logger {
    return &Logger{
        Logger: log.New(os.Stdout, "", log.LstdFlags),
        level:  level,
    }
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
    // Add context information to logs
    return l
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
    // Add structured fields to logs
    return l
}
```

Features:

- Structured logging
- Context propagation
- Log levels
- Request tracing
- Error tracking

### 2. Common Models

```go
// pkg/models/common.go
package models

type BaseResponse struct {
    Status  string      `json:"status"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

type Pagination struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
    Total    int `json:"total"`
}
```

Features:

- Standard response formats
- Common data structures
- Validation rules
- Type definitions

### 3. Configuration Package

```go
// pkg/config/config.go
package config

type BaseConfig struct {
    Environment string
    LogLevel    string
    ServiceName string
    Version     string
}

func LoadBaseConfig() *BaseConfig {
    return &BaseConfig{
        Environment: getEnvOrDefault("ENV", "development"),
        LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
        ServiceName: getEnvOrDefault("SERVICE_NAME", "unknown"),
        Version:     getEnvOrDefault("VERSION", "0.0.1"),
    }
}
```

Features:

- Environment-based configuration
- Default values
- Validation
- Type safety

### 4. Health Check Package

```go
// pkg/health/checker.go
package health

type HealthChecker struct {
    checks map[string]Check
}

type Check func() error

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]Check),
    }
}

func (h *HealthChecker) AddCheck(name string, check Check) {
    h.checks[name] = check
}
```

Features:

- Standard health check interface
- Dependency checking
- Status reporting
- Metrics collection

## Best Practices

### 1. Early Logging Implementation

- Implement logging from the start of development
- Use structured logging for better analysis
- Include context in all log messages
- Define log levels appropriately
- Add request tracing IDs

Example:

```go
func (h *Handler) CreateProfile(c *gin.Context) {
    logger := logging.FromContext(c)
    logger.Info("Creating new profile",
        "user_id", c.GetString("user_id"),
        "request_id", c.GetString("request_id"))
    // ... implementation
}
```

### 2. Synchronized Models

- Define shared model definitions
- Use code generation for model synchronization
- Implement validation at model level
- Keep models in sync across services
- Document model changes

Example:

```go
// pkg/models/profile.go
type Profile struct {
    ID        string    `json:"id" validate:"required"`
    UserID    string    `json:"user_id" validate:"required"`
    FirstName string    `json:"first_name" validate:"required"`
    LastName  string    `json:"last_name" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 3. Error Handling

- Define common error types
- Implement consistent error responses
- Add error context
- Log errors appropriately
- Handle errors at appropriate levels

Example:

```go
// pkg/errors/errors.go
type ServiceError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

func NewServiceError(code string, message string) *ServiceError {
    return &ServiceError{
        Code:    code,
        Message: message,
        Details: make(map[string]interface{}),
    }
}
```

### 4. Testing Utilities

- Shared test helpers
- Common test fixtures
- Mock implementations
- Test utilities
- Integration test helpers

Example:

```go
// pkg/testing/helpers.go
func NewTestServer(t *testing.T) *httptest.Server {
    // Create test server with common configuration
}

func NewTestClient(t *testing.T) *http.Client {
    // Create test client with common configuration
}
```

## Implementation Guidelines

1. **Service Template**

   - Use the service template for new services
   - Include all shared components
   - Follow the established patterns
   - Document deviations

2. **Code Generation**

   - Use code generation for models
   - Generate API clients
   - Generate documentation
   - Generate tests

3. **Documentation**

   - Document shared components
   - Keep documentation up to date
   - Include examples
   - Document best practices

4. **Testing**
   - Write tests for shared components
   - Include integration tests
   - Test error scenarios
   - Test performance

## Benefits

1. **Consistency**

   - Uniform logging
   - Standard error handling
   - Common response formats
   - Consistent configuration

2. **Maintainability**

   - Reduced code duplication
   - Centralized updates
   - Easier debugging
   - Better documentation

3. **Development Speed**

   - Faster service creation
   - Reusable components
   - Standard patterns
   - Less boilerplate

4. **Operational Efficiency**
   - Better monitoring
   - Easier troubleshooting
   - Consistent metrics
   - Standard health checks
