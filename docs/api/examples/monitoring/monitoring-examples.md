# Monitoring API Examples

This guide provides examples of how to use the Monitoring API for the Profile Service Microservices architecture.

## Health Checks

### Check Service Health

```bash
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"
```

Response:

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-03-20T10:00:00Z",
  "details": {
    "database": "healthy",
    "cache": "healthy",
    "event_bus": "healthy"
  }
}
```

### Go Client Example

```go
package monitoring

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type HealthClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type HealthStatus struct {
    Status    string                 `json:"status"`
    Version   string                 `json:"version"`
    Timestamp string                 `json:"timestamp"`
    Details   map[string]string      `json:"details"`
}

func NewHealthClient(baseURL, apiKey string) *HealthClient {
    return &HealthClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *HealthClient) CheckHealth() (*HealthStatus, error) {
    req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("health check failed: %s", resp.Status)
    }

    var status HealthStatus
    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        return nil, err
    }

    return &status, nil
}
```

## Metrics

### List Available Metrics

```bash
curl -X GET https://monitoring.profileservice.com/v1/metrics \
  -H "X-API-Key: your-api-key"
```

Response:

```json
[
  "profile_service.request_count",
  "profile_service.request_latency",
  "profile_service.error_count",
  "profile_service.active_connections",
  "profile_service.cache_hits",
  "profile_service.cache_misses"
]
```

### Query Metrics

```bash
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service.request_count",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m",
    "labels": {
      "service": "profile-api",
      "endpoint": "/profiles"
    }
  }'
```

Response:

```json
[
  {
    "name": "profile_service.request_count",
    "value": 150,
    "timestamp": "2024-03-20T09:00:00Z",
    "labels": {
      "service": "profile-api",
      "endpoint": "/profiles"
    }
  },
  {
    "name": "profile_service.request_count",
    "value": 175,
    "timestamp": "2024-03-20T09:05:00Z",
    "labels": {
      "service": "profile-api",
      "endpoint": "/profiles"
    }
  }
]
```

### Go Client Example

```go
package monitoring

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type MetricClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type MetricQuery struct {
    Name      string            `json:"name"`
    StartTime string            `json:"start_time"`
    EndTime   string            `json:"end_time"`
    Interval  string            `json:"interval,omitempty"`
    Labels    map[string]string `json:"labels,omitempty"`
}

type Metric struct {
    Name      string            `json:"name"`
    Value     float64           `json:"value"`
    Timestamp string            `json:"timestamp"`
    Labels    map[string]string `json:"labels,omitempty"`
}

func NewMetricClient(baseURL, apiKey string) *MetricClient {
    return &MetricClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *MetricClient) QueryMetrics(query *MetricQuery) ([]Metric, error) {
    jsonData, err := json.Marshal(query)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.baseURL+"/metrics", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("query failed: %s", resp.Status)
    }

    var metrics []Metric
    if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
        return nil, err
    }

    return metrics, nil
}
```

## Logs

### Query Logs

```bash
curl -X POST https://monitoring.profileservice.com/v1/logs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "level": "error",
    "service": "profile-api",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "failed to process request"
  }'
```

Response:

```json
[
  {
    "level": "error",
    "message": "failed to process request: database connection timeout",
    "timestamp": "2024-03-20T09:15:00Z",
    "service": "profile-api",
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "fields": {
      "request_id": "req-123",
      "endpoint": "/profiles",
      "error": "connection timeout"
    }
  }
]
```

### Go Client Example

```go
package monitoring

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type LogClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type LogQuery struct {
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
    Level     string `json:"level,omitempty"`
    Service   string `json:"service,omitempty"`
    TraceID   string `json:"trace_id,omitempty"`
    Message   string `json:"message,omitempty"`
}

type LogEntry struct {
    Level     string            `json:"level"`
    Message   string            `json:"message"`
    Timestamp string            `json:"timestamp"`
    Service   string            `json:"service"`
    TraceID   string            `json:"trace_id,omitempty"`
    Fields    map[string]string `json:"fields,omitempty"`
}

func NewLogClient(baseURL, apiKey string) *LogClient {
    return &LogClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *LogClient) QueryLogs(query *LogQuery) ([]LogEntry, error) {
    jsonData, err := json.Marshal(query)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.baseURL+"/logs", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("query failed: %s", resp.Status)
    }

    var entries []LogEntry
    if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
        return nil, err
    }

    return entries, nil
}
```

## Alerts

### List Active Alerts

```bash
curl -X GET https://monitoring.profileservice.com/v1/alerts \
  -H "X-API-Key: your-api-key"
```

Response:

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "high_error_rate",
    "severity": "critical",
    "status": "firing",
    "description": "Error rate exceeded threshold",
    "created_at": "2024-03-20T09:15:00Z",
    "updated_at": "2024-03-20T09:15:00Z"
  }
]
```

### Go Client Example

```go
package monitoring

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type AlertClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type Alert struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Severity    string    `json:"severity"`
    Status      string    `json:"status"`
    Description string    `json:"description"`
    CreatedAt   string    `json:"created_at"`
    UpdatedAt   string    `json:"updated_at"`
}

func NewAlertClient(baseURL, apiKey string) *AlertClient {
    return &AlertClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *AlertClient) ListAlerts() ([]Alert, error) {
    req, err := http.NewRequest("GET", c.baseURL+"/alerts", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to list alerts: %s", resp.Status)
    }

    var alerts []Alert
    if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
        return nil, err
    }

    return alerts, nil
}
```

## Traces

### Get Trace Information

```bash
curl -X GET "https://monitoring.profileservice.com/v1/traces?trace_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "X-API-Key: your-api-key"
```

Response:

```json
{
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "spans": [
    {
      "span_id": "span-1",
      "name": "profile-service.request",
      "start_time": "2024-03-20T09:15:00Z",
      "end_time": "2024-03-20T09:15:01Z",
      "attributes": {
        "http.method": "GET",
        "http.path": "/profiles",
        "http.status_code": "200"
      }
    },
    {
      "span_id": "span-2",
      "name": "database.query",
      "start_time": "2024-03-20T09:15:00.1Z",
      "end_time": "2024-03-20T09:15:00.5Z",
      "attributes": {
        "db.type": "postgres",
        "db.statement": "SELECT * FROM profiles"
      }
    }
  ]
}
```

### Go Client Example

```go
package monitoring

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type TraceClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

type Trace struct {
    TraceID string `json:"trace_id"`
    Spans   []Span `json:"spans"`
}

type Span struct {
    SpanID    string            `json:"span_id"`
    Name      string            `json:"name"`
    StartTime string            `json:"start_time"`
    EndTime   string            `json:"end_time"`
    Attributes map[string]string `json:"attributes"`
}

func NewTraceClient(baseURL, apiKey string) *TraceClient {
    return &TraceClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        apiKey:     apiKey,
    }
}

func (c *TraceClient) GetTrace(traceID string) (*Trace, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/traces?trace_id=%s", c.baseURL, traceID), nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get trace: %s", resp.Status)
    }

    var trace Trace
    if err := json.NewDecoder(resp.Body).Decode(&trace); err != nil {
        return nil, err
    }

    return &trace, nil
}
```
