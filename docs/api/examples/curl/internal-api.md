# Internal API cURL Examples

This document provides cURL examples for interacting with the Internal API.

## Authentication

### Service-to-Service Authentication

```bash
curl -X POST https://internal.profileservice.com/v1/auth/service-token \
  -H "Content-Type: application/json" \
  -H "X-Service-Name: profile-service" \
  -H "X-Service-Secret: your-service-secret" \
  -d '{
    "service_id": "profile-service",
    "service_secret": "your-service-secret"
  }'
```

## Service Operations

### Health Check

```bash
curl -X GET https://internal.profileservice.com/v1/health \
  -H "X-Service-Token: your-service-token"
```

### Service Status

```bash
curl -X GET https://internal.profileservice.com/v1/services/status \
  -H "X-Service-Token: your-service-token"
```

## Data Operations

### Batch Profile Creation

```bash
curl -X POST https://internal.profileservice.com/v1/profiles/batch \
  -H "X-Service-Token: your-service-token" \
  -H "Content-Type: application/json" \
  -d '{
    "profiles": [
      {
        "name": "John Doe",
        "email": "john.doe@example.com",
        "bio": "Software Engineer"
      },
      {
        "name": "Jane Smith",
        "email": "jane.smith@example.com",
        "bio": "Product Manager"
      }
    ]
  }'
```

### Bulk Profile Update

```bash
curl -X PUT https://internal.profileservice.com/v1/profiles/bulk \
  -H "X-Service-Token: your-service-token" \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {
        "id": "123",
        "name": "John Doe Updated"
      },
      {
        "id": "456",
        "bio": "Senior Product Manager"
      }
    ]
  }'
```

### Batch Profile Deletion

```bash
curl -X DELETE https://internal.profileservice.com/v1/profiles/batch \
  -H "X-Service-Token: your-service-token" \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["123", "456", "789"]
  }'
```

## Event Operations

### Publish Event

```bash
curl -X POST https://internal.profileservice.com/v1/events/publish \
  -H "X-Service-Token: your-service-token" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "profile.updated",
    "data": {
      "profile_id": "123",
      "changes": {
        "name": "John Doe Updated",
        "bio": "Senior Software Engineer"
      }
    }
  }'
```

### Subscribe to Events

```bash
curl -X POST https://internal.profileservice.com/v1/events/subscribe \
  -H "X-Service-Token: your-service-token" \
  -H "Content-Type: application/json" \
  -d '{
    "event_types": ["profile.created", "profile.updated", "profile.deleted"],
    "callback_url": "https://your-service.com/webhooks/profile-events"
  }'
```

## Error Handling

### Service Authentication Error

```bash
# Response when service token is invalid
curl -X GET https://internal.profileservice.com/v1/health \
  -H "X-Service-Token: invalid-token"
# Response: 401 Unauthorized
```

### Rate Limit Error

```bash
# Response when rate limit is exceeded
curl -X GET https://internal.profileservice.com/v1/services/status \
  -H "X-Service-Token: your-service-token"
# Response: 429 Too Many Requests
```

## Headers

### Required Headers

```bash
curl -X GET https://internal.profileservice.com/v1/health \
  -H "X-Service-Token: your-service-token" \
  -H "X-Request-ID: your-request-id" \
  -H "Content-Type: application/json"
```

### Optional Headers

```bash
curl -X GET https://internal.profileservice.com/v1/health \
  -H "X-Service-Token: your-service-token" \
  -H "X-Service-Version: 1.0.0" \
  -H "X-Environment: production"
```

## Response Examples

### Success Response

```json
{
  "data": {
    "status": "healthy",
    "services": {
      "profile-service": "up",
      "auth-service": "up",
      "event-service": "up"
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
    "code": "service_error",
    "message": "Service authentication failed",
    "details": {
      "service": "profile-service",
      "reason": "Invalid service token"
    }
  },
  "meta": {
    "request_id": "req-123"
  }
}
```

## Best Practices

1. **Service Authentication**

   - Use secure service tokens
   - Rotate service secrets regularly
   - Implement token refresh logic
   - Monitor authentication failures

2. **Error Handling**

   - Implement retry logic
   - Handle rate limiting
   - Log service errors
   - Monitor error rates

3. **Performance**

   - Use batch operations
   - Implement caching
   - Monitor response times
   - Handle timeouts

4. **Security**

   - Use internal network
   - Validate service tokens
   - Encrypt sensitive data
   - Monitor access patterns

5. **Monitoring**
   - Track service health
   - Monitor event processing
   - Log important operations
   - Set up alerts
