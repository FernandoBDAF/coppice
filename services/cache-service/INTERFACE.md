# Cache Service Interface Details

## API Endpoints

### Cache Operations

1. **Get Value**

   ```http
   GET /v1/cache/{key}
   Authorization: Bearer access_token
   ```

2. **Set Value**

   ```http
   POST /v1/cache
   Content-Type: application/json
   Authorization: Bearer access_token

   {
     "key": "user:123",
     "value": "user_data",
     "ttl": 3600
   }
   ```

3. **Delete Value**

   ```http
   DELETE /v1/cache/{key}
   Authorization: Bearer access_token
   ```

4. **Batch Operations**

   ```http
   POST /v1/cache/batch
   Content-Type: application/json
   Authorization: Bearer access_token

   {
     "operations": [
       {
         "type": "get",
         "key": "user:123"
       },
       {
         "type": "set",
         "key": "user:456",
         "value": "user_data",
         "ttl": 3600
       }
     ]
   }
   ```

### Lock Operations

1. **Acquire Lock**

   ```http
   POST /v1/locks
   Content-Type: application/json
   Authorization: Bearer access_token

   {
     "key": "resource:123",
     "ttl": 30
   }
   ```

2. **Release Lock**

   ```http
   DELETE /v1/locks/{key}
   Authorization: Bearer access_token
   ```

### Cache Management

1. **Get Stats**

   ```http
   GET /v1/stats
   Authorization: Bearer access_token
   ```

2. **Clear Cache**

   ```http
   POST /v1/cache/clear
   Content-Type: application/json
   Authorization: Bearer access_token

   {
     "pattern": "user:*"
   }
   ```

## Service Dependencies

### External Services

1. **Redis**

   - Purpose: Primary caching backend
   - Operations: GET, SET, DEL, SCAN
   - Keys: `cache:{key}`, `lock:{key}`

2. **Monitoring Service**

   - Purpose: Metrics and health monitoring
   - Integration: Prometheus metrics
   - Events: Health checks

3. **Logging Service**
   - Purpose: Operation logs
   - Integration: Structured logging
   - Events: Cache operations

### Internal Services

1. **Auth Service**

   - Purpose: Authentication and authorization
   - Integration: JWT validation
   - Events: Token validation

2. **Profile Service**

   - Purpose: User profile caching
   - Integration: Cache operations
   - Events: Profile updates

3. **Storage Service**

   - Purpose: File metadata caching
   - Integration: Cache operations
   - Events: File operations

4. **Worker Service**
   - Purpose: Job result caching
   - Integration: Cache operations
   - Events: Job completion

## Message Queue Topics

1. **Cache Events**

   ```
   cache.events.key.expired
   cache.events.key.deleted
   cache.events.pattern.cleared
   ```

2. **Lock Events**
   ```
   cache.events.lock.acquired
   cache.events.lock.released
   cache.events.lock.expired
   ```

## Response Formats

### Success Response

```json
{
  "status": "success",
  "data": {
    "key": "user:123",
    "value": "user_data",
    "ttl": 3600,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Error Response

```json
{
  "status": "error",
  "error": {
    "type": "CACHE_MISS",
    "message": "Key not found",
    "details": ["Key: user:123"]
  }
}
```

## Rate Limiting

1. **Cache Operations**

   - 1000 requests per minute per IP
   - 100 requests per second per key

2. **Lock Operations**
   - 100 requests per minute per IP
   - 10 requests per second per key

## Security Headers

1. **Required Headers**

   ```
   Authorization: Bearer <token>
   Content-Type: application/json
   ```

2. **Optional Headers**
   ```
   X-Request-ID: <uuid>
   X-Client-Version: <version>
   ```

## CORS Configuration

```go
config := cors.Config{
    AllowedOrigins: []string{
        "https://app.example.com",
        "https://api.example.com",
    },
    AllowedMethods: []string{
        "GET",
        "POST",
        "DELETE",
        "OPTIONS",
    },
    AllowedHeaders: []string{
        "Authorization",
        "Content-Type",
        "X-Request-ID",
    },
    MaxAge: 86400,
}
```
