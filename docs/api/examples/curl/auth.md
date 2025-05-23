# Auth API cURL Examples

This document provides cURL examples for interacting with the Auth API.

## Authentication

### Get Access Token

```bash
curl -X POST https://auth.profileservice.com/v1/token \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "client_credentials",
    "client_id": "your-client-id",
    "client_secret": "your-client-secret"
  }'
```

### Refresh Token

```bash
curl -X POST https://auth.profileservice.com/v1/token/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "refresh_token",
    "refresh_token": "your-refresh-token"
  }'
```

## User Management

### Register User

```bash
curl -X POST https://auth.profileservice.com/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure-password",
    "name": "John Doe"
  }'
```

### Login User

```bash
curl -X POST https://auth.profileservice.com/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure-password"
  }'
```

### Get User Profile

```bash
curl -X GET https://auth.profileservice.com/v1/users/me \
  -H "Authorization: Bearer your-access-token"
```

### Update User Profile

```bash
curl -X PUT https://auth.profileservice.com/v1/users/me \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.doe@example.com"
  }'
```

### Change Password

```bash
curl -X POST https://auth.profileservice.com/v1/users/change-password \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "current-password",
    "new_password": "new-secure-password"
  }'
```

## Role Management

### Assign Role

```bash
curl -X POST https://auth.profileservice.com/v1/roles/assign \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "role": "admin"
  }'
```

### Remove Role

```bash
curl -X DELETE https://auth.profileservice.com/v1/roles/remove \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "role": "admin"
  }'
```

### List User Roles

```bash
curl -X GET https://auth.profileservice.com/v1/roles/user/user-123 \
  -H "Authorization: Bearer your-access-token"
```

## Session Management

### List Active Sessions

```bash
curl -X GET https://auth.profileservice.com/v1/sessions \
  -H "Authorization: Bearer your-access-token"
```

### Revoke Session

```bash
curl -X DELETE https://auth.profileservice.com/v1/sessions/session-123 \
  -H "Authorization: Bearer your-access-token"
```

### Revoke All Sessions

```bash
curl -X DELETE https://auth.profileservice.com/v1/sessions/all \
  -H "Authorization: Bearer your-access-token"
```

## Error Handling

### Invalid Credentials

```bash
# Response when credentials are invalid
curl -X POST https://auth.profileservice.com/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "wrong-password"
  }'
# Response: 401 Unauthorized
```

### Token Expired

```bash
# Response when token is expired
curl -X GET https://auth.profileservice.com/v1/users/me \
  -H "Authorization: Bearer expired-token"
# Response: 401 Unauthorized
```

### Rate Limit Exceeded

```bash
# Response when rate limit is exceeded
curl -X POST https://auth.profileservice.com/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure-password"
  }'
# Response: 429 Too Many Requests
```

## Headers

### Required Headers

```bash
curl -X GET https://auth.profileservice.com/v1/users/me \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json"
```

### Optional Headers

```bash
curl -X GET https://auth.profileservice.com/v1/users/me \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en-US" \
  -H "X-Client-Version: 1.0.0"
```

## Response Examples

### Success Response

```json
{
  "data": {
    "user": {
      "id": "user-123",
      "email": "user@example.com",
      "name": "John Doe",
      "roles": ["user", "admin"],
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z"
    }
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
    "code": "invalid_credentials",
    "message": "Invalid email or password",
    "details": {
      "attempts_remaining": 4,
      "lockout_time": "2024-03-20T10:05:00Z"
    }
  },
  "meta": {
    "request_id": "req-123"
  }
}
```

## Best Practices

1. **Authentication**

   - Use secure password storage
   - Implement token refresh
   - Handle session management
   - Monitor login attempts

2. **Security**

   - Use HTTPS only
   - Implement rate limiting
   - Validate all inputs
   - Monitor security events

3. **Error Handling**

   - Implement retry logic
   - Handle token expiration
   - Log security events
   - Monitor error rates

4. **Performance**

   - Cache user data
   - Optimize token validation
   - Monitor response times
   - Handle high load

5. **Monitoring**
   - Track authentication attempts
   - Monitor session usage
   - Log security events
   - Set up alerts
