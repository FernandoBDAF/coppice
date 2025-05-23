# Monitoring Troubleshooting Guide

This guide provides solutions for common monitoring issues in the Profile Service Microservices architecture.

## 1. Health Check Issues

### Service Unhealthy

**Symptoms:**

- Health check returns "unhealthy" status
- Service is not responding
- High error rates

**Troubleshooting Steps:**

1. Check service logs for errors
2. Verify database connectivity
3. Check cache connection
4. Verify event bus status
5. Check resource usage

**Example Commands:**

```bash
# Check service logs
curl -X POST https://monitoring.profileservice.com/v1/logs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "level": "error",
    "service": "profile-api"
  }'

# Check resource usage
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_resource_usage",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z"
  }'
```

## 2. High Error Rates

### API Errors

**Symptoms:**

- Increased error responses
- Failed requests
- Timeout errors

**Troubleshooting Steps:**

1. Check error logs
2. Verify request patterns
3. Check dependent services
4. Review rate limits
5. Check resource constraints

**Example Commands:**

```bash
# Check error logs
curl -X POST https://monitoring.profileservice.com/v1/logs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "level": "error",
    "service": "profile-api",
    "message": "failed to process request"
  }'

# Check error rates
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_error_count",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 3. Performance Issues

### High Latency

**Symptoms:**

- Slow response times
- Increased request duration
- Timeout errors

**Troubleshooting Steps:**

1. Check response time metrics
2. Review database performance
3. Check cache hit rates
4. Verify network latency
5. Check resource usage

**Example Commands:**

```bash
# Check response times
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_request_duration_seconds",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'

# Check cache performance
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_cache_hits",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 4. Event System Issues

### Event Processing Delays

**Symptoms:**

- Events not processed
- Increased queue size
- Processing delays

**Troubleshooting Steps:**

1. Check event queue size
2. Review processing rates
3. Check dead letter queue
4. Verify event consumer status
5. Check event bus health

**Example Commands:**

```bash
# Check queue size
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_event_queue_size",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'

# Check processing rates
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_events_processed_total",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 5. Resource Issues

### High Resource Usage

**Symptoms:**

- High CPU usage
- Memory pressure
- Disk space issues
- Network congestion

**Troubleshooting Steps:**

1. Check resource metrics
2. Review process stats
3. Check connection pools
4. Verify disk usage
5. Check network stats

**Example Commands:**

```bash
# Check CPU usage
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "process_cpu_seconds_total",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'

# Check memory usage
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "process_resident_memory_bytes",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 6. Security Issues

### Authentication Failures

**Symptoms:**

- Increased auth failures
- API key issues
- Rate limit hits

**Troubleshooting Steps:**

1. Check auth logs
2. Review API key usage
3. Check rate limits
4. Verify security config
5. Review access patterns

**Example Commands:**

```bash
# Check auth failures
curl -X POST https://monitoring.profileservice.com/v1/logs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "level": "error",
    "service": "profile-api",
    "message": "authentication failed"
  }'

# Check rate limits
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_rate_limit_hits_total",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 7. Monitoring System Issues

### Data Collection Problems

**Symptoms:**

- Missing metrics
- Delayed data
- Incomplete data
- Collection errors

**Troubleshooting Steps:**

1. Check collector status
2. Verify data sources
3. Review collection config
4. Check storage system
5. Verify network connectivity

**Example Commands:**

```bash
# Check collector status
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"

# Check data collection
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_metrics_collected_total",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z",
    "interval": "5m"
  }'
```

## 8. Common Solutions

### Service Restart

```bash
# Check service status
curl -X GET https://monitoring.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"

# Restart service (if needed)
kubectl rollout restart deployment/profile-service
```

### Clear Cache

```bash
# Check cache status
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_cache_status",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z"
  }'

# Clear cache (if needed)
curl -X POST https://profile-service:8080/v1/internal/cache/clear \
  -H "X-API-Key: your-api-key"
```

### Scale Resources

```bash
# Check resource usage
curl -X POST https://monitoring.profileservice.com/v1/metrics \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "profile_service_resource_usage",
    "start_time": "2024-03-20T09:00:00Z",
    "end_time": "2024-03-20T10:00:00Z"
  }'

# Scale service (if needed)
kubectl scale deployment/profile-service --replicas=3
```

## 9. Prevention

### Regular Maintenance

1. Review metrics daily
2. Check alert thresholds
3. Monitor trends
4. Update documentation

### Proactive Monitoring

1. Set up predictive alerts
2. Monitor capacity trends
3. Review performance patterns
4. Check security events

### Documentation

1. Update runbooks
2. Document solutions
3. Share knowledge
4. Review procedures
