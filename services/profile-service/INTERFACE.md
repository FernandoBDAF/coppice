# Profile Service Interface Documentation

## API Overview

The Profile Service provides a RESTful API for managing user profiles. All endpoints are prefixed with `/api/v1/profiles`.

## Authentication

All endpoints require authentication using a Bearer token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### Create Profile

```http
POST /api/v1/profiles
```

Creates a new user profile.

#### Request Body

```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "metadata": {
    "preferences": {
      "theme": "dark",
      "notifications": true
    }
  }
}
```

#### Response

```json
{
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2024-02-20T10:00:00Z",
    "metadata": {
      "version": 1,
      "preferences": {
        "theme": "dark",
        "notifications": true
      }
    }
  }
}
```

### Get Profile

```http
GET /api/v1/profiles/{id}
```

Retrieves a user profile by ID.

#### Response

```json
{
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2024-02-20T10:00:00Z",
    "metadata": {
      "version": 1,
      "preferences": {
        "theme": "dark",
        "notifications": true
      }
    }
  }
}
```

### Update Profile

```http
PUT /api/v1/profiles/{id}
```

Updates an existing user profile.

#### Request Body

```json
{
  "name": "John Updated",
  "metadata": {
    "preferences": {
      "theme": "light",
      "notifications": false
    }
  }
}
```

#### Response

```json
{
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Updated",
    "created_at": "2024-02-20T10:00:00Z",
    "updated_at": "2024-02-20T10:00:00Z",
    "metadata": {
      "version": 2,
      "preferences": {
        "theme": "light",
        "notifications": false
      }
    }
  }
}
```

### Delete Profile

```http
DELETE /api/v1/profiles/{id}
```

Soft deletes a user profile.

#### Response

```json
{
  "status": "success",
  "message": "Profile deleted successfully"
}
```

## Error Responses

All endpoints may return the following error responses:

### Validation Error

```json
{
  "type": "VALIDATION_ERROR",
  "message": "Invalid request data",
  "details": ["Email is required", "Name must be at least 2 characters"]
}
```

### Authentication Error

```json
{
  "type": "AUTHENTICATION_ERROR",
  "message": "Invalid or expired token"
}
```

### Not Found Error

```json
{
  "type": "NOT_FOUND_ERROR",
  "message": "Profile not found"
}
```

### Service Unavailable Error

```json
{
  "type": "SERVICE_UNAVAILABLE_ERROR",
  "message": "Service temporarily unavailable"
}
```

## Rate Limiting

The API implements rate limiting:

- 100 requests per minute per IP
- 1000 requests per hour per IP

Rate limit headers are included in all responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1613822400
```

## Metrics

The service exposes Prometheus metrics at `/metrics`:

```
# HELP profile_service_requests_total Total number of requests
# TYPE profile_service_requests_total counter
profile_service_requests_total{endpoint="/api/v1/profiles",method="POST"} 100

# HELP profile_service_request_duration_seconds Request duration in seconds
# TYPE profile_service_request_duration_seconds histogram
profile_service_request_duration_seconds_bucket{endpoint="/api/v1/profiles",method="POST",le="0.1"} 90
```

## Health Check

```http
GET /health
```

Response:

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-02-20T10:00:00Z"
}
```
