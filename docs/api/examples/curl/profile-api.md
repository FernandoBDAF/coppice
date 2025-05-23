# Profile API cURL Examples

This document provides cURL examples for interacting with the Profile API.

## Authentication

### Get API Token

```bash
curl -X POST https://api.profileservice.com/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "your-client-id",
    "client_secret": "your-client-secret"
  }'
```

### Using API Key

```bash
curl -X GET https://api.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"
```

## Profile Operations

### List Profiles

```bash
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
```

### Get Profile by ID

```bash
curl -X GET https://api.profileservice.com/v1/profiles/{id} \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
```

### Create Profile

```bash
curl -X POST https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "bio": "Software Engineer",
    "image_urls": ["https://example.com/image1.jpg"]
  }'
```

### Update Profile

```bash
curl -X PUT https://api.profileservice.com/v1/profiles/{id} \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "bio": "Senior Software Engineer"
  }'
```

### Delete Profile

```bash
curl -X DELETE https://api.profileservice.com/v1/profiles/{id} \
  -H "Authorization: Bearer your-token"
```

## Error Handling

### Rate Limit Example

```bash
# Response when rate limit is exceeded
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
# Response: 429 Too Many Requests
```

### Validation Error Example

```bash
# Response when validation fails
curl -X POST https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "email": "invalid-email"
  }'
# Response: 400 Bad Request
```

### Authentication Error Example

```bash
# Response when token is invalid
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer invalid-token" \
  -H "Content-Type: application/json"
# Response: 401 Unauthorized
```

## Query Parameters

### Pagination

```bash
curl -X GET "https://api.profileservice.com/v1/profiles?page=1&limit=10" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
```

### Filtering

```bash
curl -X GET "https://api.profileservice.com/v1/profiles?name=John&email=example.com" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
```

### Sorting

```bash
curl -X GET "https://api.profileservice.com/v1/profiles?sort=name&order=desc" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json"
```

## Headers

### Required Headers

```bash
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "X-Request-ID: your-request-id"
```

### Optional Headers

```bash
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en-US" \
  -H "X-Client-Version: 1.0.0"
```

## Response Examples

### Success Response

```json
{
  "data": {
    "id": "123",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "bio": "Software Engineer",
    "image_urls": ["https://example.com/image1.jpg"],
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
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
    "code": "validation_error",
    "message": "Invalid input data",
    "details": {
      "name": "Name is required",
      "email": "Invalid email format"
    }
  },
  "meta": {
    "request_id": "req-123"
  }
}
```

## Best Practices

1. **Error Handling**

   - Always check response status codes
   - Parse error responses for details
   - Implement retry logic for transient errors

2. **Authentication**

   - Store tokens securely
   - Implement token refresh logic
   - Use API keys for service-to-service communication

3. **Rate Limiting**

   - Implement exponential backoff
   - Monitor rate limit headers
   - Cache responses when appropriate

4. **Security**

   - Use HTTPS only
   - Validate all inputs
   - Sanitize sensitive data in logs

5. **Performance**
   - Use compression when available
   - Implement connection pooling
   - Cache responses appropriately
