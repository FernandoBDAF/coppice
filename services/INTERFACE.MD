# Service Interface Documentation

## Interface Standards

### REST API Standards

| Standard        | Description                                          | Implementation           |
| --------------- | ---------------------------------------------------- | ------------------------ |
| OpenAPI/Swagger | All REST APIs must be documented using OpenAPI 3.0   | See `/docs/api/openapi/` |
| Authentication  | JWT-based authentication with Clerk integration      | See `auth-api` service   |
| Error Handling  | Standardized error responses with codes and messages | See `pkg/errors`         |
| Rate Limiting   | Token bucket algorithm with Redis                    | See `pkg/ratelimit`      |

### gRPC Standards

| Standard         | Description                                    | Implementation        |
| ---------------- | ---------------------------------------------- | --------------------- |
| Protocol Buffers | All gRPC services must use Protocol Buffers v3 | See `/proto/`         |
| Service Patterns | Unary and streaming patterns defined           | See `pkg/grpc`        |
| Error Handling   | Standardized gRPC error codes and messages     | See `pkg/grpc/errors` |

### Message Queue Standards

| Standard       | Description                           | Implementation         |
| -------------- | ------------------------------------- | ---------------------- |
| Message Format | JSON with schema validation           | See `pkg/queue`        |
| Event Patterns | Event sourcing and CQRS patterns      | See `pkg/events`       |
| Error Handling | Dead letter queues and retry policies | See `pkg/queue/errors` |

## Current Service Interactions

### External Interfaces

#### Profile API Service

| Endpoint                | Method | Purpose        | Authentication |
| ----------------------- | ------ | -------------- | -------------- |
| `/api/v1/profiles`      | GET    | List profiles  | Required       |
| `/api/v1/profiles/{id}` | GET    | Get profile    | Required       |
| `/api/v1/profiles`      | POST   | Create profile | Required       |
| `/api/v1/profiles/{id}` | PUT    | Update profile | Required       |
| `/api/v1/profiles/{id}` | DELETE | Delete profile | Required       |

#### Auth API Service

| Endpoint                | Method | Purpose            | Authentication |
| ----------------------- | ------ | ------------------ | -------------- |
| `/api/v1/auth/login`    | POST   | User login         | None           |
| `/api/v1/auth/register` | POST   | User registration  | None           |
| `/api/v1/auth/verify`   | POST   | Token verification | Required       |

### Internal Interfaces

#### Profile Storage Service

| Service          | Method          | Purpose          | Protocol |
| ---------------- | --------------- | ---------------- | -------- |
| `ProfileService` | `GetProfile`    | Get profile data | gRPC     |
| `ProfileService` | `CreateProfile` | Create profile   | gRPC     |
| `ProfileService` | `UpdateProfile` | Update profile   | gRPC     |
| `ProfileService` | `DeleteProfile` | Delete profile   | gRPC     |

#### Cache Service

| Service        | Method   | Purpose            | Protocol |
| -------------- | -------- | ------------------ | -------- |
| `CacheService` | `Get`    | Get cached data    | gRPC     |
| `CacheService` | `Set`    | Set cached data    | gRPC     |
| `CacheService` | `Delete` | Delete cached data | gRPC     |

### Queue Interactions

#### Profile Queue Service

| Queue             | Purpose           | Message Format | Consumer       |
| ----------------- | ----------------- | -------------- | -------------- |
| `profile-updates` | Profile updates   | JSON           | Worker Service |
| `profile-deletes` | Profile deletions | JSON           | Worker Service |

## Cross-Service Models

### Shared Models

Location: `common/models/`

| Model     | Purpose              | Used By               |
| --------- | -------------------- | --------------------- |
| `Profile` | User profile data    | Profile API, Storage  |
| `User`    | User authentication  | Auth API, Profile API |
| `Token`   | Authentication token | Auth API, Profile API |

### Model Compatibility

- All services must use the same model versions
- Model changes require coordination across services
- Versioning strategy: Semantic versioning

## Planned Interactions

### Future Implementations

1. **Cache Service Integration**

   - Reference: `docs/architecture/cache-service.md`
   - Dependencies: Redis, Profile API
   - Status: Planned

2. **Queue Service Integration**

   - Reference: `docs/architecture/queue-service.md`
   - Dependencies: RabbitMQ, Worker Service
   - Status: Planned

3. **Monitoring Service Integration**
   - Reference: `docs/architecture/monitoring-service.md`
   - Dependencies: Prometheus, Grafana
   - Status: Planned

## Notes

- All interfaces must be documented in OpenAPI/Swagger or Protocol Buffers
- Changes to interfaces require coordination across services
- Interface versioning follows semantic versioning
- All services must implement health check endpoints
- All services must implement metrics endpoints
- All services must implement logging endpoints
