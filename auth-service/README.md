# Auth Service

JWT-based authentication microservice (Node.js / TypeScript / Express), part of the
`microservices` monorepo. Issues and validates JWTs for `api-service`, manages users,
and enforces rate limiting + account lockout.

## Features

- User authentication (login, logout, token refresh)
- JWT access + refresh tokens (`jsonwebtoken` v9)
- Token validation endpoint consumed by `api-service`
- User management (CRUD, activation, roles)
- Rate limiting and account lockout
- Structured logging (pino), Prometheus metrics (prom-client), health/readiness probes

## Integration with api-service

`api-service` validates tokens by calling (see `CONTRACTS.md` §3 — this shape is frozen):

```
POST /v1/auth/token/validate
Authorization: Bearer <token>
Content-Type: application/json

{ "token": "<jwt>" }
```

200 response:

```json
{
  "status": "success",
  "message": "Token is valid",
  "data": { "valid": true, "user": { "id": "...", "email": "...", "role": "..." } }
}
```

An invalid/expired token returns `401` (api-service treats any non-200 response as invalid).

## API Endpoints

### Authentication (`/v1/auth`)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/auth/login` | POST | User login, returns access + refresh tokens |
| `/v1/auth/token/validate` | POST | Validate JWT (used by api-service) |
| `/v1/auth/token/refresh` | POST | Exchange a refresh token for a new token pair |
| `/v1/auth/logout` | POST | Logout (invalidates the token for audit purposes) |

### User Management (`/v1/users`)
| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/users` | POST | none | Create user (self-registration) |
| `/v1/users` | GET | admin | List users (paginated) |
| `/v1/users/me` | GET | any user | Get current user profile |
| `/v1/users/:id` | GET/PUT/DELETE | admin | User CRUD |
| `/v1/users/email/:email` | GET | admin | Lookup user by email |
| `/v1/users/:id/activate` | PATCH | admin | Activate user |
| `/v1/users/:id/deactivate` | PATCH | admin | Deactivate user |
| `/v1/users/:id/role` | PATCH | admin | Change user role |

### Health & Monitoring
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Liveness + dependency status (DB) |
| `/ready` | GET | Readiness probe (real PostgreSQL check) |
| `/live` | GET | Bare process liveness (no dependency checks) |
| `/metrics` | GET | Prometheus metrics (`prom-client`, default + HTTP histograms) |

In non-production environments, interactive API docs are served at `/api-docs`
(OpenAPI JSON at `/api-docs.json`).

## Environment Variables

Names marked **frozen** are pinned in `../CONTRACTS.md` §4 and must not be renamed.

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `NODE_ENV` **(frozen)** | development | No | `development` \| `production` \| `test` |
| `PORT` **(frozen)** | 3000 | No | HTTP port (compose sets `3000`) |
| `DATABASE_HOST` **(frozen)** | postgres | No | PostgreSQL host |
| `DATABASE_PORT` **(frozen)** | 5432 | No | PostgreSQL port |
| `DATABASE_NAME` **(frozen)** | auth_db | No | Database name |
| `DATABASE_USER` **(frozen)** | auth_user | No | Database user |
| `DATABASE_PASSWORD` **(frozen)** | - | **Yes** | Database password |
| `DATABASE_POOL_MAX` | 20 | No | Max PG pool connections |
| `DATABASE_SSL` | false | No | Enable TLS to PostgreSQL |
| `JWT_SECRET` **(frozen)** | - | **Yes** | JWT signing secret (min 32 chars) |
| `JWT_ACCESS_TOKEN_EXPIRY` | 15m | No | Access token lifetime |
| `JWT_REFRESH_TOKEN_EXPIRY` | 7d | No | Refresh token lifetime |
| `RATE_LIMIT_WINDOW_MS` | 900000 | No | Rate limit window (15 min) |
| `RATE_LIMIT_MAX_REQUESTS` | 100 | No | Max auth requests per window per IP |
| `ACCOUNT_LOCKOUT_ATTEMPTS` | 5 | No | Failed attempts before lockout |
| `ACCOUNT_LOCKOUT_DURATION_MS` | 1800000 | No | Lockout duration (30 min) |
| `PASSWORD_MIN_LENGTH` | 8 | No | Minimum password length |
| `METRICS_ENABLED` | true | No | Enable `/metrics` + Prometheus collection |
| `METRICS_PREFIX` | auth_service_ | No | Prefix for all emitted metric names |
| `LOG_LEVEL` **(frozen)** | info | No | pino log level |
| `LOG_PRETTY` **(frozen)** | false | No | Human-readable logs instead of JSON |
| `API_SERVICE_URL` | http://api-service:8080 | No | Reserved for future use; not currently called |

## Local Development

```bash
# From auth-service/
npm install

# Requires a reachable PostgreSQL (see Quickstart below for the compose one),
# plus DATABASE_PASSWORD and JWT_SECRET (min 32 chars) in the environment.
npm run dev          # tsx watch, applies migrations on boot

npm run typecheck
npm run lint
npm run build         # emits ESM to dist/
npm start              # node dist/server.js

npm test
npm run test:coverage
```

## Quickstart (whole stack via root docker-compose.yml)

This service is built and orchestrated from the **repository root** compose file, not
standalone:

```bash
# From the repository root (not auth-service/)
docker compose up -d postgres
docker compose up -d auth-migrate   # one-shot: applies migrations/*.sql
docker compose up -d auth-service

# or just bring up everything:
docker compose up -d
```

`auth-service` listens on `http://localhost:3000` (container port 3000, per
`CONTRACTS.md` §1). Compose supplies `DATABASE_HOST=postgres`,
`DATABASE_NAME=auth_db`, `DATABASE_USER=auth_user`, `DATABASE_PASSWORD=auth_password`,
and a dev `JWT_SECRET` — see the root `docker-compose.yml` for exact values.

Note: migrations run twice in the compose path — once via the one-shot `auth-migrate`
container (raw `psql`) and again via this service's own startup migration runner. Both
are idempotent (`CREATE ... IF NOT EXISTS`, constraint-existence checks), so this is safe.

## Database

PostgreSQL schema (`migrations/001_create_users_table.sql`, idempotent):

- `users` — accounts: `id`, `email`, `hashed_password`, `salt`, `role`, `is_active`,
  `failed_attempts`, `locked_until`, timestamps. Passwords are hashed with bcrypt over
  `password + salt`, where `salt` is a per-user random value stored alongside the hash.
- `auth_audit_logs` — authentication event log, FK to `users(id)`.

Migrations in `migrations/` are applied automatically on service startup
(`src/infrastructure/database/migrations.ts`), tracked in a `migrations` table with a
checksum per file.
