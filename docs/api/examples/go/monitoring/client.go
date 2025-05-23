package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents a Monitoring API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new Monitoring API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// HealthStatus represents the health status of services
type HealthStatus struct {
	Status    string                 `json:"status"`
	Services  map[string]ServiceInfo `json:"services"`
	Timestamp time.Time              `json:"timestamp"`
}

// ServiceInfo represents information about a service
type ServiceInfo struct {
	Status    string  `json:"status"`
	Latency   string  `json:"latency"`
	ErrorRate float64 `json:"error_rate"`
}

// GetHealth retrieves the health status of all services
func (c *Client) GetHealth(ctx context.Context) (*HealthStatus, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var health HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &health, nil
}

// Metric represents a metric value
type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

// GetMetrics retrieves metrics for a service
func (c *Client) GetMetrics(ctx context.Context, service, metric string, interval string) ([]Metric, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/metrics", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("service", service)
	q.Add("metric", metric)
	q.Add("interval", interval)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metrics []Metric
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return metrics, nil
}

// SubmitMetric submits a custom metric
func (c *Client) SubmitMetric(ctx context.Context, metric Metric) error {
	jsonData, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/metrics/custom", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// LogEntry represents a log entry
type LogEntry struct {
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Service   string            `json:"service"`
	Labels    map[string]string `json:"labels"`
}

// QueryLogs retrieves logs based on query parameters
func (c *Client) QueryLogs(ctx context.Context, service, level string, start, end time.Time) ([]LogEntry, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/logs", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("service", service)
	q.Add("level", level)
	q.Add("start", start.Format(time.RFC3339))
	q.Add("end", end.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var logs []LogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return logs, nil
}

// SearchLogs searches logs using a query string
func (c *Client) SearchLogs(ctx context.Context, query string, start, end time.Time, limit int) ([]LogEntry, error) {
	payload := map[string]interface{}{
		"query":      query,
		"start_time": start.Format(time.RFC3339),
		"end_time":   end.Format(time.RFC3339),
		"limit":      limit,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/logs/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var logs []LogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return logs, nil
}

// Alert represents an alert configuration
type Alert struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Condition     string              `json:"condition"`
	Duration      string              `json:"duration"`
	Severity      string              `json:"severity"`
	Notifications map[string][]string `json:"notifications"`
}

// ListAlerts retrieves all active alerts
func (c *Client) ListAlerts(ctx context.Context) ([]Alert, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/alerts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var alerts []Alert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return alerts, nil
}

// CreateAlert creates a new alert
func (c *Client) CreateAlert(ctx context.Context, alert Alert) (*Alert, error) {
	jsonData, err := json.Marshal(alert)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/alerts", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdAlert Alert
	if err := json.NewDecoder(resp.Body).Decode(&createdAlert); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdAlert, nil
}

// UpdateAlert updates an existing alert
func (c *Client) UpdateAlert(ctx context.Context, alertID string, updates map[string]interface{}) (*Alert, error) {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/v1/alerts/"+alertID, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updatedAlert Alert
	if err := json.NewDecoder(resp.Body).Decode(&updatedAlert); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedAlert, nil
}

// DeleteAlert deletes an alert
func (c *Client) DeleteAlert(ctx context.Context, alertID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/v1/alerts/"+alertID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Trace represents a distributed trace
type Trace struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Operation string    `json:"operation"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`
	Spans     []Span    `json:"spans"`
}

// Span represents a span in a trace
type Span struct {
	ID        string            `json:"id"`
	TraceID   string            `json:"trace_id"`
	ParentID  string            `json:"parent_id"`
	Name      string            `json:"name"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  string            `json:"duration"`
	Tags      map[string]string `json:"tags"`
}

// GetTrace retrieves a specific trace
func (c *Client) GetTrace(ctx context.Context, traceID string) (*Trace, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/traces/"+traceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var trace Trace
	if err := json.NewDecoder(resp.Body).Decode(&trace); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &trace, nil
}

// SearchTraces searches for traces based on criteria
func (c *Client) SearchTraces(ctx context.Context, service, operation string, start, end time.Time, minDuration string) ([]Trace, error) {
	payload := map[string]interface{}{
		"service":      service,
		"operation":    operation,
		"start_time":   start.Format(time.RFC3339),
		"end_time":     end.Format(time.RFC3339),
		"min_duration": minDuration,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/traces/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var traces []Trace
	if err := json.NewDecoder(resp.Body).Decode(&traces); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return traces, nil
}
