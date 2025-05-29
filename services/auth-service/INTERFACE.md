# Auth Service Interface Details

## API Endpoints

### Authentication Endpoints

1. **User Registration**

   ```http
   POST /v1/auth/register
   Content-Type: application/json

   {
     "email": "user@example.com",
     "password": "securePassword123",
     "role": "user"
   }
   ```

2. **User Login**

   ```http
   POST /v1/auth/login
   Content-Type: application/json

   {
     "email": "user@example.com",
     "password": "securePassword123"
   }
   ```

3. **Token Refresh**

   ```http
   POST /v1/auth/token/refresh
   Content-Type: application/json

   {
     "refresh_token": "refresh_token"
   }
   ```

4. **Token Validation**

   ```http
   POST /v1/auth/token/validate
   Content-Type: application/json

   {
     "token": "access_token"
   }
   ```

### User Management Endpoints

1. **Get Current User**

   ```http
   GET /v1/users/me
   Authorization: Bearer access_token
   ```

2. **Get User by ID**
   ```http
   GET /v1/users/{id}
   Authorization: Bearer access_token
   ```

### OAuth Endpoints

1. **OAuth Authorization**

   ```http
   GET /v1/oauth/authorize
   ```

2. **OAuth Token**

   ```http
   POST /v1/oauth/token
   Content-Type: application/json

   {
     "grant_type": "authorization_code",
     "code": "authorization_code"
   }
   ```

3. **OAuth User Info**
   ```http
   GET /v1/oauth/userinfo
   Authorization: Bearer oauth_access_token
   ```

## Service Dependencies

### External Services

1. **Clerk Authentication**

   - Purpose: User authentication and management
   - Integration: REST API and Webhooks
   - Events: User creation, updates, deletion

2. **Redis**

   - Purpose: Session storage and token caching
   - Operations: GET, SET, DEL
   - Keys: `session:{id}`, `token:{id}`

3. **PostgreSQL**
   - Purpose: User data and role storage
   - Tables: users, roles, permissions
   - Operations: CRUD operations

### Internal Services

1. **Profile Service**

   - Purpose: User profile management
   - Integration: REST API
   - Events: Profile updates

2. **Cache Service**

   - Purpose: Distributed caching
   - Integration: Redis protocol
   - Operations: Cache invalidation

3. **Monitoring Service**

   - Purpose: Metrics and health monitoring
   - Integration: Prometheus metrics
   - Events: Health checks

4. **Worker Service**
   - Purpose: Background tasks
   - Integration: Message queue
   - Tasks: Email notifications, cleanup

## Message Queue Topics

1. **Authentication Events**

   ```
   auth.events.user.created
   auth.events.user.updated
   auth.events.user.deleted
   auth.events.session.created
   auth.events.session.expired
   ```

2. **Notification Events**
   ```
   auth.notifications.email.verification
   auth.notifications.password.reset
   auth.notifications.security.alert
   ```

## Webhook Endpoints

1. **Clerk Webhooks**

   ```http
   POST /v1/webhooks/clerk
   Content-Type: application/json

   {
     "type": "user.created",
     "data": {
       "id": "user_id",
       "email": "user@example.com"
     }
   }
   ```

2. **OAuth Webhooks**

   ```http
   POST /v1/webhooks/oauth
   Content-Type: application/json

   {
     "type": "token.refreshed",
     "data": {
       "user_id": "user_id",
       "token": "new_token"
     }
   }
   ```

## Response Formats

### Success Response

```json
{
  "status": "success",
  "message": "Operation successful",
  "data": {
    // Response data specific to the endpoint
  }
}
```

### Error Response

```json
{
  "status": "error",
  "message": "Error message",
  "errors": [
    {
      "field": "field_name",
      "message": "Error details"
    }
  ]
}
```

## Rate Limiting

1. **Authentication Endpoints**

   - 20 requests per minute per IP
   - 100 requests per minute per user

2. **API Endpoints**
   - 100 requests per minute per token
   - 1000 requests per minute per IP

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
        "PUT",
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
