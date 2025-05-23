# Event System Documentation

This guide explains how to use the event system in the Profile Service Microservices architecture.

## Overview

The event system enables asynchronous communication between services using an event-driven architecture. Events are published to an event bus and consumed by interested services.

## Event Types

### 1. Profile Events

```json
{
  "type": "profile.created",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "bio": "Software Engineer",
    "created_at": "2024-03-20T10:00:00Z"
  },
  "timestamp": "2024-03-20T10:00:00Z"
}
```

```json
{
  "type": "profile.updated",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "bio": "Senior Software Engineer",
    "updated_at": "2024-03-20T11:00:00Z"
  },
  "timestamp": "2024-03-20T11:00:00Z"
}
```

```json
{
  "type": "profile.deleted",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "deleted_at": "2024-03-20T12:00:00Z"
  },
  "timestamp": "2024-03-20T12:00:00Z"
}
```

### 2. Task Events

```json
{
  "type": "task.created",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "profile_cleanup",
    "status": "pending",
    "created_at": "2024-03-20T10:00:00Z"
  },
  "timestamp": "2024-03-20T10:00:00Z"
}
```

```json
{
  "type": "task.completed",
  "payload": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "profile_cleanup",
    "status": "completed",
    "completed_at": "2024-03-20T10:05:00Z"
  },
  "timestamp": "2024-03-20T10:05:00Z"
}
```

## Publishing Events

### Using the Internal API

```bash
curl -X POST http://profile-service:8080/v1/internal/events \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "type": "profile.created",
    "payload": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "bio": "Software Engineer",
      "created_at": "2024-03-20T10:00:00Z"
    },
    "timestamp": "2024-03-20T10:00:00Z"
  }'
```

### Using Go Client

```go
package events

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type EventClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type Event struct {
    Type      string      `json:"type"`
    Payload   interface{} `json:"payload"`
    Timestamp string      `json:"timestamp"`
}

func NewEventClient(baseURL, apiKey string) *EventClient {
    return &EventClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *EventClient) PublishEvent(event *Event) error {
    jsonData, err := json.Marshal(event)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", c.baseURL+"/internal/events", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return fmt.Errorf("failed to publish event: %s", resp.Status)
    }

    return nil
}
```

## Consuming Events

### Using Go Consumer

```go
package events

import (
    "context"
    "encoding/json"
    "log"
)

type EventConsumer struct {
    eventBus EventBus
    handlers map[string]EventHandler
}

type EventHandler func(ctx context.Context, payload interface{}) error

func NewEventConsumer(eventBus EventBus) *EventConsumer {
    return &EventConsumer{
        eventBus:  eventBus,
        handlers:  make(map[string]EventHandler),
    }
}

func (c *EventConsumer) RegisterHandler(eventType string, handler EventHandler) {
    c.handlers[eventType] = handler
}

func (c *EventConsumer) Start(ctx context.Context) error {
    events := c.eventBus.Subscribe(ctx)

    for event := range events {
        if handler, ok := c.handlers[event.Type]; ok {
            if err := handler(ctx, event.Payload); err != nil {
                log.Printf("Error handling event %s: %v", event.Type, err)
            }
        }
    }

    return nil
}
```

## Event Processing Patterns

### 1. At-Least-Once Delivery

```go
func (c *EventConsumer) processWithRetry(ctx context.Context, event *Event, handler EventHandler) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := handler(ctx, event.Payload)
        if err == nil {
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Second * time.Duration(i+1)):
            continue
        }
    }
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

### 2. Dead Letter Queue

```go
func (c *EventConsumer) processWithDLQ(ctx context.Context, event *Event, handler EventHandler) error {
    err := handler(ctx, event.Payload)
    if err != nil {
        // Send to dead letter queue
        dlqEvent := &Event{
            Type:      "dlq." + event.Type,
            Payload:   event,
            Timestamp: time.Now().UTC().Format(time.RFC3339),
        }
        return c.eventBus.Publish(ctx, dlqEvent)
    }
    return nil
}
```

## Best Practices

1. **Event Naming**

   - Use dot notation for event types (e.g., `profile.created`)
   - Use past tense for completed actions
   - Be specific about the event type

2. **Event Payload**

   - Include all necessary data for consumers
   - Keep payloads small and focused
   - Use consistent data structures

3. **Error Handling**

   - Implement retry mechanisms
   - Use dead letter queues for failed events
   - Log event processing errors

4. **Monitoring**
   - Track event processing latency
   - Monitor event queue sizes
   - Alert on processing failures

## Next Steps

1. 🔄 Implement event versioning
2. 🔄 Add event schema validation
3. 🔄 Create event replay functionality
4. 🔄 Add event correlation IDs
5. 🔄 Implement event ordering
6. 🔄 Add event compression
