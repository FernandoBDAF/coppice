# Storage Service Technical Context

## Internal Architecture

### Core Components

1. **Storage Layer** (`internal/storage/`)

   - Database operations
   - Data persistence
   - Transaction management
   - Connection pooling
   - Query optimization
   - Data validation

2. **Service Layer** (`internal/service/`)

   - Business logic
   - Data transformation
   - Integration with shared libraries
   - Error handling
   - Request validation
   - Response formatting

3. **API Layer** (`internal/api/`)

   - REST API endpoints
   - gRPC service
   - Health checks
   - Metrics endpoints
   - Request validation
   - Response formatting

4. **Integration Layer** (`internal/integration/`)
   - Shared libraries integration
   - API services communication
   - Circuit breaking
   - Retry mechanisms
   - Service discovery
   - Load balancing

### Design Patterns

1. **Repository Pattern**

   - Data access abstraction
   - CRUD operations
   - Query optimization
   - Transaction management

2. **Factory Pattern**

   - Connection factory
   - Client factory
   - Repository factory
   - Service factory

3. **Strategy Pattern**

   - Query strategies
   - Caching strategies
   - Retry strategies
   - Error handling strategies

4. **Observer Pattern**
   - Database monitoring
   - Performance tracking
   - Health monitoring
   - Error tracking

### Frameworks and Libraries

1. **Database Framework**

   - PostgreSQL driver
   - Connection pooling
   - Query builder
   - Migration tool

2. **Web Framework**

   - Gin for HTTP routing
   - gRPC for RPC
   - Validator for request validation
   - JWT-Go for authentication

3. **Testing**

   - Go testing package
   - Testify for assertions
   - Mockery for mocking
   - Testcontainers for integration tests

4. **Utilities**
   - Zap for logging
   - Viper for configuration
   - Wire for dependency injection
   - UUID for IDs

### Data Models

1. **Profile Model**

```go
type Profile struct {
    ID          string    `json:"id"`
    Email       string    `json:"email"`
    Name        string    `json:"name"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   time.Time `json:"deleted_at,omitempty"`
    Metadata    Metadata  `json:"metadata"`
}

type Metadata struct {
    Version     int       `json:"version"`
    LastLogin   time.Time `json:"last_login,omitempty"`
    Preferences map[string]interface{} `json:"preferences"`
}
```

2. **Address Model**

```go
type Address struct {
    ID          string    `json:"id"`
    ProfileID   string    `json:"profile_id"`
    Type        string    `json:"type"`
    Street      string    `json:"street"`
    City        string    `json:"city"`
    State       string    `json:"state"`
    Country     string    `json:"country"`
    PostalCode  string    `json:"postal_code"`
    IsDefault   bool      `json:"is_default"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

3. **Contact Model**

```go
type Contact struct {
    ID          string    `json:"id"`
    ProfileID   string    `json:"profile_id"`
    Type        string    `json:"type"`
    Value       string    `json:"value"`
    IsVerified  bool      `json:"is_verified"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Storage Strategy

1. **Data Types**

   - Profile data
   - Address data
   - Contact data
   - Metadata

2. **Storage Patterns**

   - CRUD operations
   - Batch operations
   - Query optimization
   - Transaction management

3. **Persistence Strategy**
   - Data persistence
   - Soft deletion
   - Version control
   - Audit logging

### Error Handling

1. **Error Types**

```go
type ErrorType string

const (
    ErrNotFound        ErrorType = "NOT_FOUND_ERROR"
    ErrValidation      ErrorType = "VALIDATION_ERROR"
    ErrDatabase        ErrorType = "DATABASE_ERROR"
    ErrDuplicateEntry ErrorType = "DUPLICATE_ENTRY_ERROR"
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
   - Profile ID
   - Operation type
   - Duration
   - Error details

### Metrics Collection

1. **Database Metrics**

   - Query performance
   - Connection pool
   - Transaction rates
   - Error rates

2. **API Metrics**
   - Request rates
   - Response times
   - Error rates
   - Resource usage

### Security Implementation

1. **Authentication**

   - JWT validation
   - API key validation
   - Role-based access

2. **Authorization**

   - Data access control
   - Operation permissions
   - Resource limits

3. **Data Security**
   - Data encryption
   - Secure connections
   - Access logging
   - Audit trail

### Testing Strategy

1. **Unit Tests**

   - Repository tests
   - Service tests
   - API tests
   - Validation tests

2. **Integration Tests**

   - Database integration
   - API integration
   - Service integration
   - End-to-end tests

3. **Performance Tests**
   - Query performance
   - API performance
   - Load testing
   - Stress testing
