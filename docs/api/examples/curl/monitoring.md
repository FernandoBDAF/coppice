# Monitoring API cURL Examples

This document provides cURL examples for interacting with the Monitoring API.

## Health Checks

### Service Health

```bash
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"
```

### Detailed Health Check

```bash
curl -X GET https://monitoring.profileservice.com/v1/health/detailed \
  -H "X-API-Key: your-api-key"
```

### Health Check History

```bash
curl -X GET "https://monitoring.profileservice.com/v1/health/history?start=2024-03-19T00:00:00Z&end=2024-03-20T00:00:00Z" \
  -H "X-API-Key: your-api-key"
```

## Metrics

### Get Service Metrics

```bash
curl -X GET "https://monitoring.profileservice.com/v1/metrics?service=profile-service&metric=request_count&interval=1h" \
  -H "X-API-Key: your-api-key"
```

### Get Custom Metrics

```bash
curl -X GET "https://monitoring.profileservice.com/v1/metrics/custom?name=custom_metric&start=2024-03-19T00:00:00Z&end=2024-03-20T00:00:00Z" \
  -H "X-API-Key: your-api-key"
```

### Submit Custom Metric

```bash
curl -X POST https://monitoring.profileservice.com/v1/metrics/custom \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "custom_metric",
    "value": 42,
    "labels": {
      "service": "profile-service",
      "environment": "production"
    }
  }'
```

## Logs

### Query Logs

```bash
curl -X GET "https://monitoring.profileservice.com/v1/logs?service=profile-service&level=error&start=2024-03-19T00:00:00Z&end=2024-03-20T00:00:00Z" \
  -H "X-API-Key: your-api-key"
```

### Search Logs

```bash
curl -X POST https://monitoring.profileservice.com/v1/logs/search \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "error AND service:profile-service",
    "start_time": "2024-03-19T00:00:00Z",
    "end_time": "2024-03-20T00:00:00Z",
    "limit": 100
  }'
```

### Get Log Patterns

```bash
curl -X GET "https://monitoring.profileservice.com/v1/logs/patterns?service=profile-service&start=2024-03-19T00:00:00Z&end=2024-03-20T00:00:00Z" \
  -H "X-API-Key: your-api-key"
```

## Alerts

### List Active Alerts

```bash
curl -X GET https://monitoring.profileservice.com/v1/alerts \
  -H "X-API-Key: your-api-key"
```

### Create Alert

```bash
curl -X POST https://monitoring.profileservice.com/v1/alerts \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "high_error_rate",
    "condition": "error_rate > 0.1",
    "duration": "5m",
    "severity": "critical",
    "notifications": {
      "email": ["team@example.com"],
      "slack": ["#alerts"]
    }
  }'
```

### Update Alert

```bash
curl -X PUT https://monitoring.profileservice.com/v1/alerts/alert-123 \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "condition": "error_rate > 0.05",
    "severity": "warning"
  }'
```

### Delete Alert

```bash
curl -X DELETE https://monitoring.profileservice.com/v1/alerts/alert-123 \
  -H "X-API-Key: your-api-key"
```

## Traces

### Get Trace

```bash
curl -X GET https://monitoring.profileservice.com/v1/traces/trace-123 \
  -H "X-API-Key: your-api-key"
```

### Search Traces

```bash
curl -X POST https://monitoring.profileservice.com/v1/traces/search \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "service": "profile-service",
    "operation": "create_profile",
    "start_time": "2024-03-19T00:00:00Z",
    "end_time": "2024-03-20T00:00:00Z",
    "min_duration": "100ms"
  }'
```

### Get Trace Statistics

```bash
curl -X GET "https://monitoring.profileservice.com/v1/traces/stats?service=profile-service&start=2024-03-19T00:00:00Z&end=2024-03-20T00:00:00Z" \
  -H "X-API-Key: your-api-key"
```

## Error Handling

### Invalid API Key

```bash
# Response when API key is invalid
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: invalid-key"
# Response: 401 Unauthorized
```

### Rate Limit Exceeded

```bash
# Response when rate limit is exceeded
curl -X GET https://monitoring.profileservice.com/v1/metrics \
  -H "X-API-Key: your-api-key"
# Response: 429 Too Many Requests
```

## Headers

### Required Headers

```bash
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json"
```

### Optional Headers

```bash
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -H "X-Environment: production" \
  -H "X-Client-Version: 1.0.0"
```

## Response Examples

### Success Response

```json
{
  "data": {
    "status": "healthy",
    "services": {
      "profile-service": {
        "status": "up",
        "latency": "50ms",
        "error_rate": "0.01"
      },
      "auth-service": {
        "status": "up",
        "latency": "30ms",
        "error_rate": "0.005"
      }
    },
    "timestamp": "2024-03-20T10:00:00Z"
  },
  "meta": {
    "request_id": "req-123"
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Too many requests",
    "details": {
      "limit": 100,
      "remaining": 0,
      "reset_time": "2024-03-20T10:01:00Z"
    }
  },
  "meta": {
    "request_id": "req-123"
  }
}
```

## Best Practices

1. **Monitoring**

   - Set up comprehensive alerts
   - Monitor key metrics
   - Track error rates
   - Monitor performance

2. **Logging**

   - Use appropriate log levels
   - Include context in logs
   - Implement log rotation
   - Monitor log patterns

3. **Metrics**

   - Define clear metrics
   - Set up dashboards
   - Monitor trends
   - Set up alerts

4. **Security**

   - Use API keys
   - Implement rate limiting
   - Monitor access
   - Audit logs

5. **Performance**
   - Optimize queries
   - Cache responses
   - Monitor latency
   - Handle high load
