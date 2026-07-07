# API Service

Consolidated Go API service that provides Profile CRUD, document upload, and
task submission with direct access to PostgreSQL, Redis, RabbitMQ, and MinIO.
Go 1.24, Gin, sqlx + lib/pq.

## Features
- Profile CRUD API with a cache-aside Redis layer (single entries and list
  pages are both cached and invalidated on writes)
- Document upload to MinIO with metadata in PostgreSQL and a `document.process`
  task published for the graphrag-service to pick up
- RabbitMQ task publishing (direct, durable exchanges with per-task DLQs)
- Token validation via auth-service, guarded by a circuit breaker
- Health/readiness checks and Prometheus metrics on a dedicated port
- Structured JSON logging (Zap), graceful shutdown on SIGTERM/SIGINT

## Endpoints

Main API (port `8080`):
```
GET    /health                                  liveness
GET    /ready                                   readiness (checks Postgres + Redis + RabbitMQ)

GET    /api/v1/profiles
POST   /api/v1/profiles
GET    /api/v1/profiles/:id
PUT    /api/v1/profiles/:id
DELETE /api/v1/profiles/:id

POST   /api/v1/profiles/:id/tasks               generic task (routing_key + type + payload in body)
POST   /api/v1/profiles/:id/tasks/email
POST   /api/v1/profiles/:id/tasks/image
POST   /api/v1/profiles/:id/tasks/profile
GET    /api/v1/profiles/:id/documents            paginated: ?page=&page_size=

POST   /api/v1/documents/upload                  multipart/form-data: file, profile_id
GET    /api/v1/documents/:id
GET    /api/v1/documents/:id/status
GET    /api/v1/documents/:id/download            returns a 15-minute presigned MinIO URL
DELETE /api/v1/documents/:id
```
All `/api/v1/*` routes require `Authorization: Bearer <jwt>`, validated against auth-service.

Metrics server (own port `8081`):
```
GET    /metrics                                  Prometheus, separate port per platform contract
```

## Configuration

Config is loaded via Viper (optional `config.yaml` in the working directory or
`./configs`, overridden by the environment variables below; unset values fall
back to the defaults shown).

| Variable | Default | Notes |
|---|---|---|
| `API_POSTGRES_DSN` | `postgres://postgres:postgres@localhost:5432/api_db?sslmode=disable` | lib/pq DSN |
| `API_REDIS_HOST` | `localhost` | |
| `API_REDIS_PORT` | `6379` | |
| `API_REDIS_PASSWORD` | _(empty)_ | |
| `API_RABBITMQ_HOSTS` | `localhost:5672` | single `host:port` or comma-separated list |
| `API_RABBITMQ_PASSWORD` | `guest` | |
| `API_AUTH_URL` | `http://auth-service:3000` | POST `{AUTH_URL}/v1/auth/token/validate` |
| `MINIO_ENDPOINT` | `minio:9000` | |
| `MINIO_ACCESS_KEY` | _(empty, required if endpoint set)_ | |
| `MINIO_SECRET_KEY` | _(empty, required if endpoint set)_ | |
| `MINIO_USE_SSL` | `false` | |
| `MINIO_BUCKET_NAME` | `documents-raw` | auto-created at startup if missing |

The main HTTP port (`8080`) and metrics port (`8081`) are fixed by the
platform contract and are not environment-overridable.

## Local Run

From the repo root, start infrastructure:
```
docker compose up -d postgres redis rabbitmq minio
```

Apply database migrations (one-shot container, applies `api-service/migrations/*.up.sql`):
```
docker compose up api-migrate
```

Run the service (from `api-service/`). The Postgres/Redis/RabbitMQ defaults
above already target `localhost`, matching the ports docker-compose publishes
to the host; auth-service calls need `API_AUTH_URL` pointed at wherever
auth-service is reachable from your host (e.g. `http://localhost:3000` if it's
also running via compose with its port published):
```
API_AUTH_URL=http://localhost:3000 \
MINIO_ACCESS_KEY=minioadmin MINIO_SECRET_KEY=minioadmin \
make run
```

To run the whole stack in containers instead, use `docker compose up -d api-service`
from the repo root (builds this service's Dockerfile and wires it to the
compose network, where the frozen defaults for `API_AUTH_URL` etc. apply directly).

## Tests
```
make test
```
Covers profile cache-invalidation behavior, task envelope construction, and the
auth-client circuit breaker (invalid tokens must not trip it; 5xx responses must).
