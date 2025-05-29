# Logging Base Library

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

The Logging Base Library provides a standardized foundation for implementing structured logging across all microservices. It offers consistent log formatting, context propagation, and log enrichment while allowing services to extend with their specific logging needs.

## Architecture

### Core Components

1. **Logger**

   ```go
   // File: logger.go
   type Logger struct {
       serviceName string
       level       Level
       fields      map[string]interface{}
       output      io.Writer
       formatter   Formatter
       hooks       []Hook
   }

   type Level int

   const (
       DebugLevel Level = iota
       InfoLevel
       WarnLevel
       ErrorLevel
       FatalLevel
   )
   ```

2. **Formatter**

   ```go
   // File: formatter.go
   type Formatter interface {
       Format(entry *Entry) ([]byte, error)
   }

   type JSONFormatter struct {
       TimestampFormat string
       PrettyPrint     bool
   }

   type Entry struct {
       Timestamp time.Time
       Level     Level
       Message   string
       Fields    map[string]interface{}
       Service   string
       TraceID   string
       SpanID    string
   }
   ```

3. **Context**

   ```go
   // File: context.go
   type Context struct {
       TraceID    string
       SpanID     string
       UserID     string
       RequestID  string
       Fields     map[string]interface{}
   }

   func WithContext(ctx context.Context) *Context {
       return &Context{
           TraceID:   trace.SpanContextFromContext(ctx).TraceID().String(),
           SpanID:    trace.SpanContextFromContext(ctx).SpanID().String(),
           RequestID: ctx.Value("request_id").(string),
           Fields:    make(map[string]interface{}),
       }
   }
   ```

4. **Hook**

   ```go
   // File: hook.go
   type Hook interface {
       Levels() []Level
       Fire(*Entry) error
   }

   type ErrorHook struct {
       errorHandler func(error)
   }

   type MetricsHook struct {
       collector *monitoring.BaseCollector
   }
   ```

## Service Integration

### Basic Integration

```go
// File: services/profile-api/internal/logging/logger.go
package logging

import (
    "github.com/your-org/logging"
)

type ProfileLogger struct {
    *logging.Logger
    customFields map[string]interface{}
}

func NewProfileLogger() *ProfileLogger {
    base := logging.NewLogger("profile-api")
    return &ProfileLogger{
        Logger:       base,
        customFields: make(map[string]interface{}),
    }
}

func (l *ProfileLogger) Initialize() error {
    // Set default fields
    l.SetDefaultFields()

    // Add custom hooks
    l.AddHooks()

    return nil
}
```

### Custom Fields Implementation

```go
// File: services/profile-api/internal/logging/fields.go
package logging

func (l *ProfileLogger) SetDefaultFields() {
    l.WithFields(map[string]interface{}{
        "environment": os.Getenv("ENVIRONMENT"),
        "version":     os.Getenv("VERSION"),
        "region":      os.Getenv("REGION"),
    })
}

func (l *ProfileLogger) WithProfileFields(profile *Profile) *ProfileLogger {
    return &ProfileLogger{
        Logger: l.Logger.WithFields(map[string]interface{}{
            "profile_id": profile.ID,
            "user_id":    profile.UserID,
            "type":       profile.Type,
        }),
    }
}
```

### Hook Implementation

```go
// File: services/profile-api/internal/logging/hooks.go
package logging

type ErrorMetricsHook struct {
    collector *monitoring.BaseCollector
}

func (h *ErrorMetricsHook) Levels() []logging.Level {
    return []logging.Level{
        logging.ErrorLevel,
        logging.FatalLevel,
    }
}

func (h *ErrorMetricsHook) Fire(entry *logging.Entry) error {
    h.collector.IncErrorCounter(
        entry.Service,
        entry.Fields["error_type"].(string),
        entry.Fields["error_code"].(string),
    )
    return nil
}
```

## Configuration

### Base Configuration

```yaml
# File: config/base.yaml
logging:
  base:
    level: info
    format: json
    timestamp_format: "2006-01-02T15:04:05.000Z"
    pretty_print: false

  fields:
    service: true
    environment: true
    version: true
    region: true

  hooks:
    error_metrics: true
    error_notification: true

  output:
    stdout: true
    stderr: true
```

### Service-Specific Configuration

```yaml
# File: services/profile-api/config/logging.yaml
logging:
  custom:
    level: debug

  fields:
    profile_type: true
    user_type: true
    operation: true

  hooks:
    profile_metrics: true
    validation_metrics: true

  output:
    file:
      enabled: true
      path: /var/log/profile-api.log
      max_size: 100MB
      max_backups: 3
```

## Usage Examples

### Basic Logging

```go
// File: services/profile-api/internal/handler/profile.go
package handler

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
    logger := h.logger.WithContext(r.Context())

    logger.Info("Creating new profile", map[string]interface{}{
        "user_id": r.Header.Get("User-ID"),
        "type":    r.URL.Query().Get("type"),
    })

    profile, err := h.service.CreateProfile(r.Context(), req)
    if err != nil {
        logger.Error("Failed to create profile", map[string]interface{}{
            "error": err.Error(),
            "code":  "PROFILE_CREATE_ERROR",
        })
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    logger.Info("Profile created successfully", map[string]interface{}{
        "profile_id": profile.ID,
        "duration":   time.Since(start).Seconds(),
    })
}
```

### Structured Logging

```go
// File: services/profile-api/internal/service/profile.go
package service

func (s *Service) ValidateProfile(ctx context.Context, profile *Profile) error {
    logger := s.logger.WithContext(ctx).WithProfileFields(profile)

    logger.Debug("Starting profile validation", map[string]interface{}{
        "validation_type": "full",
        "fields":         profile.Fields,
    })

    if err := s.validateFields(profile); err != nil {
        logger.Error("Field validation failed", map[string]interface{}{
            "error": err.Error(),
            "code":  "VALIDATION_ERROR",
            "field": err.Field,
        })
        return err
    }

    logger.Info("Profile validation successful", map[string]interface{}{
        "duration": time.Since(start).Seconds(),
    })
    return nil
}
```

## Best Practices

1. **Log Levels**

   - Use appropriate log levels
   - Be consistent with level usage
   - Include relevant context
   - Avoid sensitive data
   - Use structured logging

2. **Context**

   - Always include request context
   - Propagate trace IDs
   - Add relevant fields
   - Use consistent field names
   - Include timing information

3. **Error Logging**

   - Log errors with context
   - Include error codes
   - Add stack traces
   - Use error hooks
   - Track error metrics

4. **Performance**
   - Use appropriate log levels
   - Avoid expensive operations
   - Use async logging when needed
   - Monitor log volume
   - Use log sampling

## Cross-References

- [Logging Patterns](../../reference-materials/development/patterns/logging-patterns.md)
- [Error Handling Patterns](../../reference-materials/development/patterns/error-patterns.md)
- [Context Patterns](../../reference-materials/development/patterns/context-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Tracing Patterns](../../reference-materials/development/patterns/tracing-patterns.md)
