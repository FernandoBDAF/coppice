# Auth Service Technical Context

## Internal Architecture

### Core Components

1. **API Layer** (`internal/api/`)

   - REST API endpoints using Gin framework
   - Request validation using validator
   - Response formatting
   - Error handling middleware
   - Rate limiting middleware
   - CORS configuration

2. **Authentication Module** (`internal/auth/`)

   - JWT token generation and validation
   - Password hashing using bcrypt
   - Session management with Redis
   - OAuth 2.0 integration
   - Clerk integration layer

3. **Authorization Module** (`internal/auth/rbac/`)

   - Role-based access control (RBAC)
   - Permission management
   - Policy enforcement
   - Access token validation

4. **Service Layer** (`internal/service/`)

   - Business logic implementation
   - Transaction management
   - Error handling
   - Event publishing

5. **Repository Layer** (`internal/repository/`)
   - Database operations using GORM
   - Redis operations
   - Data validation
   - Query optimization

### Design Patterns

1. **Repository Pattern**

   - Abstracts database operations
   - Provides consistent data access interface
   - Enables easy testing and mocking

2. **Service Layer Pattern**

   - Implements business logic
   - Coordinates between repositories
   - Handles transactions

3. **Middleware Pattern**

   - Authentication middleware
   - Authorization middleware
   - Logging middleware
   - Error handling middleware

4. **Factory Pattern**
   - Service factory
   - Repository factory
   - Client factory

### Frameworks and Libraries

1. **Web Framework**

   - Gin for HTTP routing and middleware
   - Validator for request validation
   - JWT-Go for token handling

2. **Database**

   - GORM for PostgreSQL operations
   - Redis for session storage
   - Migrations using golang-migrate

3. **Testing**

   - Go testing package
   - Testify for assertions
   - Mockery for mocking

4. **Monitoring**
   - Prometheus for metrics
   - OpenTelemetry for tracing
   - Zap for logging

### Data Models

1. **User Model**

```go
type User struct {
    ID        string    `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string
    Role      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

2. **Session Model**

```go
type Session struct {
    ID        string    `gorm:"primaryKey"`
    UserID    string    `gorm:"index"`
    Token     string    `gorm:"uniqueIndex"`
    ExpiresAt time.Time
    CreatedAt time.Time
}
```

3. **Role Model**

```go
type Role struct {
    ID          string    `gorm:"primaryKey"`
    Name        string    `gorm:"uniqueIndex"`
    Permissions []string  `gorm:"type:text[]"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Caching Strategy

1. **Session Cache**

   - Redis-based session storage
   - TTL-based expiration
   - Distributed locking for concurrent access

2. **Token Cache**

   - Redis-based token blacklist
   - TTL-based expiration
   - Batch operations for efficiency

3. **User Cache**
   - Redis-based user data cache
   - Invalidation on updates
   - TTL-based expiration

### Error Handling

1. **Error Types**

```go
type ErrorType string

const (
    ErrInvalidInput    ErrorType = "INVALID_INPUT"
    ErrUnauthorized    ErrorType = "UNAUTHORIZED"
    ErrForbidden       ErrorType = "FORBIDDEN"
    ErrNotFound        ErrorType = "NOT_FOUND"
    ErrInternal        ErrorType = "INTERNAL_ERROR"
)
```

2. **Error Response**

```go
type ErrorResponse struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Details []string  `json:"details,omitempty"`
}
```

### Logging Strategy

1. **Structured Logging**

   - JSON format
   - Contextual fields
   - Log levels
   - Request tracing

2. **Log Fields**
   - Request ID
   - User ID
   - Action
   - Duration
   - Error details

### Metrics Collection

1. **Authentication Metrics**

   - Login attempts
   - Registration attempts
   - Token validations
   - Session creations

2. **Performance Metrics**
   - Response times
   - Error rates
   - Cache hit rates
   - Database query times

### Security Implementation

1. **Password Security**

   - Bcrypt hashing
   - Salt generation
   - Password policies
   - Rate limiting

2. **Token Security**

   - JWT signing
   - Token rotation
   - Blacklisting
   - Expiration handling

3. **Session Security**
   - Secure session storage
   - Session invalidation
   - Device tracking
   - Concurrent session limits

### Testing Strategy

1. **Unit Tests**

   - Service layer tests
   - Repository layer tests
   - Utility function tests
   - Mock dependencies

2. **Integration Tests**

   - API endpoint tests
   - Database integration tests
   - Redis integration tests
   - External service integration

3. **Performance Tests**
   - Load testing
   - Stress testing
   - Endurance testing
   - Scalability testing
