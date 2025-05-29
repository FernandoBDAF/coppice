INITIAL CONTEXT FOR LLM - never change the context-----------------------------
-> THIS SECTION IS A GUIDELINE TO THE LLM CONSIDER BEFORE WORKING IN THIS FILE, DO NOT CHANGE THIS

-> GOES OF THE README FILE:

- This file serves as the technical documentation of the service, providing a comprehensive overview of the codebase
- It should document:
  - Service architecture and design decisions
  - Component structure and relationships
  - API endpoints and interfaces
  - Dependencies and integration points
  - Configuration and deployment details
- This is the primary reference for understanding the technical implementation
- This file should be in sync with the `/TRACKER&MANAGER.md` where development progress and tasks are tracked
- While TRACKER&MANAGER.md focuses on "what" and "when", this file focuses on "how" and "why"

-> CONSIDERER BEFORE UPDATING THIS FILE:

- Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.
- The changes in this file need to be incremental or to update informations that you confidentilly have knowlegde, they should not be guesses. If there are questions or uncertanty add comments asking for clarification instead.
- Check the `../../microservices/docs` folder for comprehensive project details, including architecture, development guidelines, and integration points. This will help in making informed decisions, haver better context and updates to the development plan. Always compare the implementation of this project with the plan described in the docs and whenever there are inconsistancies, add comments.
- Consider structuring this documentation separating the components and describing each one of them and adding sections to how they interact with each other - because this will be very dinamic and updated during the development process it will make clear what to update after each change
- This documentation is focusing in the current component that is part of a larger project, to have a more sistemic view check `../../microservices/TRACKER&MANAGER` and `../../microservices/README`
- Do not forget to be LLM focus, so because this will be used
- For LLM-specific guidelines and patterns, refer to [LLM Integration Guide](../../docs/llm/README.md)

---

# Storage Service

## Overview

The Storage Service is a critical component of our microservices architecture, responsible for data persistence and database operations. It provides a robust and scalable solution for managing user profiles, addresses, and contact information while ensuring data integrity, consistency, and security.

## Role in the System

The Storage Service interacts with several components in our microservices ecosystem:

### Internal Services

- **Auth Service**: Validates authentication tokens and permissions
- **Profile Service**: Manages user profile data and operations
- **Cache Service**: Coordinates data caching strategies
- **Queue Service**: Handles asynchronous data operations
- **Worker Service**: Processes background data tasks
- **Monitoring Service**: Tracks performance and health metrics

### External Services

- **PostgreSQL**: Primary database for data persistence
- **Redis**: Caching and rate limiting

## Main Functionalities

### 1. Data Management

- Profile CRUD operations
- Address management
- Contact information handling
- Data validation and integrity checks
- Soft delete functionality
- Data versioning

### 2. Performance Optimization

- Query optimization
- Connection pooling
- Caching strategies
- Batch operations
- Index management

### 3. Security Features

- Data encryption
- Access control
- Audit logging
- Rate limiting
- Input validation

### 4. Monitoring and Health

- Health checks
- Performance metrics
- Error tracking
- Resource monitoring
- Database metrics

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL 14.0
- Redis 6.2
- Make (optional, for using Makefile commands)

### Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-org/microservices.git
   cd microservices/services/storage-service
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Configure Environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start the Service**

   ```bash
   # Using Go
   go run cmd/storage-service/main.go

   # Using Docker
   docker-compose up -d
   ```

### Configuration

Essential environment variables:

```env
# Service Configuration
SERVICE_NAME=storage-service
SERVICE_PORT=8080
SERVICE_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=profiles
DB_USER=profile_storage
DB_PASSWORD=your_password
DB_MAX_CONNECTIONS=20
DB_CONNECTION_TIMEOUT=5s

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_DB=0

# Security Configuration
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key
```

### Running with Docker

1. **Build the Image**

   ```bash
   docker build -t storage-service:latest .
   ```

2. **Run the Container**
   ```bash
   docker run -p 8080:8080 \
     --env-file .env \
     storage-service:latest
   ```

## Development

### Common Tasks

1. **Running Tests**

   ```bash
   # Unit tests
   go test ./internal/...

   # Integration tests
   go test ./tests/integration/...

   # All tests
   go test ./...
   ```

2. **Building the Service**

   ```bash
   go build -o storage-service ./cmd/storage-service
   ```

3. **Running Linter**
   ```bash
   golangci-lint run
   ```

### Project Structure

```
storage-service/
├── cmd/
│   └── storage-service/    # Service entry point
├── internal/
│   ├── api/               # API handlers
│   ├── service/           # Business logic
│   ├── storage/           # Database operations
│   └── integration/       # Service integration
├── pkg/
│   └── models/            # Data models
├── tests/
│   ├── integration/       # Integration tests
│   └── unit/             # Unit tests
├── configs/              # Configuration files
├── scripts/              # Utility scripts
└── docs/                 # Documentation
```

## Documentation

For more detailed information, refer to:

- [CONTEXT.md](./CONTEXT.md): Technical architecture and design decisions
- [INTERFACE.md](./INTERFACE.md): API documentation and service interfaces
- [TRACKER.md](./TRACKER.md): Development progress and planned features

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request

## License

MIT License

## Implementation Status

### Current State

1. **Storage Layer**

   - [x] Database operations implementation
   - [x] Data persistence
   - [x] Transaction management
   - [x] Connection pooling
   - [x] Connection health checks
   - [x] Retry mechanisms with backoff

2. **Service Layer**

   - [x] Business logic implementation
   - [x] Data transformation
   - [x] Shared libraries integration
   - [x] Error handling
   - [x] Email uniqueness validation
   - [x] Request correlation tracking

3. **Integration Layer**
   - [x] Shared libraries integration
   - [x] API services communication
   - [x] Circuit breaking
   - [x] Retry mechanisms
   - [x] Request body validation
   - [x] Connection health monitoring

### Recent Improvements

1. **Request Handling**

   - Added content type validation
   - Implemented request body buffering
   - Added request size limits (1MB)
   - Improved error handling and logging
   - Added correlation IDs for request tracking
   - Enhanced request validation

2. **Connection Management**

   - Configured connection pool settings
   - Added connection health checks
   - Implemented retry logic with exponential backoff
   - Added connection backoff strategy
   - Improved error handling for connection issues
   - Added connection metrics

3. **Data Validation**
   - Implemented email uniqueness check
   - Added request body validation
   - Enhanced error handling for validation failures
   - Added detailed error logging
   - Improved error categorization

### Implementation Details

1. **Request Processing**

```go
// Request validation and processing
func (h *ProfileHandler) createProfile(w http.ResponseWriter, r *http.Request) {
    // Validate content type
    if r.Header.Get("Content-Type") != "application/json" {
        h.sendError(w, http.StatusBadRequest, "Content-Type must be application/json", nil)
        return
    }

    // Validate content length
    if r.ContentLength == 0 {
        h.sendError(w, http.StatusBadRequest, "Empty request body", nil)
        return
    }

    // Read and buffer request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        h.sendError(w, http.StatusBadRequest, "Failed to read request body", err)
        return
    }

    // Validate body size
    const maxBodySize = 1 << 20 // 1MB
    if len(body) > maxBodySize {
        h.sendError(w, http.StatusRequestEntityTooLarge, "Request body too large", nil)
        return
    }
}
```

2. **Connection Management**

```go
// Connection pool configuration
func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
    // Configure connection pool
    db.SetMaxOpenConns(100)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(1 * time.Minute)
    return &ProfileRepository{db: db}
}

// Connection health check
func (r *ProfileRepository) checkConnectionHealth(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    var result int
    err := r.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
    if err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    return nil
}
```

3. **Email Uniqueness Validation**

```go
// Email uniqueness check
func (s *ProfileService) CreateProfile(ctx context.Context, req *models.ProfileRequest) (*models.Profile, error) {
    // Check if email is already in use
    existingProfile, err := s.repo.GetByEmail(ctx, req.Email)
    if err != nil {
        if !errors.Is(err, repository.ErrNotFound) {
            return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
        }
    }

    if existingProfile != nil {
        return nil, ErrDuplicateEmail
    }
}
```

### Error Handling

1. **Request Errors**

   - Invalid content type
   - Empty request body
   - Request body too large
   - Invalid JSON format
   - Missing required fields

2. **Connection Errors**

   - Connection timeout
   - Connection reset
   - Connection refused
   - Broken pipe
   - I/O timeout

3. **Validation Errors**
   - Duplicate email
   - Invalid email format
   - Missing required fields
   - Invalid field values

### Monitoring and Metrics

1. **Request Metrics**

   - Request duration
   - Request size
   - Error rates by type
   - Success rates

2. **Connection Metrics**

   - Connection pool size
   - Connection health status
   - Connection errors
   - Connection latency

3. **Validation Metrics**
   - Validation errors by type
   - Duplicate email attempts
   - Invalid request rates

## API Endpoints

### 1. Profile Storage

```http
GET /api/v1/storage/profiles
GET /api/v1/storage/profiles/{id}
POST /api/v1/storage/profiles
PUT /api/v1/storage/profiles/{id}
DELETE /api/v1/storage/profiles/{id}
```

### 2. Batch Operations

```http
POST /api/v1/storage/profiles/batch
PUT /api/v1/storage/profiles/batch
DELETE /api/v1/storage/profiles/batch
```

### 3. Query Operations

```http
POST /api/v1/storage/profiles/query
GET /api/v1/storage/profiles/search
```

### 4. Health and Metrics

```http
GET /health
GET /ready
GET /metrics
```

## Error Types

### 1. Database Errors

- Connection errors
- Query errors
- Transaction errors
- Constraint errors
- Timeout errors

### 2. Storage Errors

- File system errors
- Permission errors
- Space errors
- IO errors
- Lock errors

### Recovery Strategies

### 1. Database Recovery

- Connection retry
- Query retry
- Transaction rollback
- Connection pool recovery
- Error logging

### 2. Storage Recovery

- File system recovery
- Permission recovery
- Space management
- IO retry
- Lock recovery

## Cross-References

- [Storage Service Patterns](../../reference-materials/development/patterns/storage-service-patterns.md)
- [Service Integration Patterns](../../reference-materials/development/patterns/service-integration-patterns.md)
- [Monitoring Patterns](../../reference-materials/development/patterns/monitoring-patterns.md)
- [Security Patterns](../../reference-materials/development/patterns/security-patterns.md)
- [Error Handling Patterns](../../reference-materials/development/patterns/error-handling-patterns.md)
