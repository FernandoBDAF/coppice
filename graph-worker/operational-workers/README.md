# operational-workers

Three independent Go consumer binaries — `email-worker`, `image-worker`,
`profile-worker` — that each consume one RabbitMQ queue published by
`api-service`, simulate the work, and ack/nack accordingly. They share
`internal/common` (queue consumer, base worker lifecycle, HTTP health
server, metrics helpers) and have one `internal/processors/<name>` package
each.

No real SMTP/image-processing/profile-store integration is wired up yet —
processors validate the payload against
`graph-worker/shared/contracts/MESSAGE_FORMAT.md` and simulate the work
with a short delay and structured log lines.

## Build & test

```bash
go build ./...
go vet ./...
go test ./...
```

## Run locally

```bash
RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run ./cmd/email-worker
```

Each worker exposes `GET /health` and `GET /ready` on `HEALTH_PORT`
(default `8080`).

## Contract (see CONTRACTS.md §2 and shared/contracts/ROUTING_KEYS.md)

| Worker | Exchange | Queue | Routing key | Prefetch | Message TTL | DLQ |
|---|---|---|---|---|---|---|
| email-worker | `email-tasks` | `email-processing` | `email.send` | 5 | 1h | `email-tasks.dlx` / `email-processing.dlq` |
| image-worker | `image-tasks` | `image-processing` | `image.process` | 1 | 6h | `image-tasks.dlx` / `image-processing.dlq` |
| profile-worker | `tasks-exchange` * | `profile-processing` | `profile.task` | 2 | 24h * | `tasks-exchange.dlx` / `profile-processing.dlq` |

\* **Known contract/code drift, not a bug in this module.** CONTRACTS.md
and `shared/contracts/ROUTING_KEYS.md` say the profile exchange is
`profile-tasks` with a 1h TTL. api-service's actual publisher
(`api-service/internal/domain/task/model.go`, `DefaultRoutingMap["profile.task"]`)
declares and publishes to exchange `tasks-exchange` with a 24h TTL / 7d
DLQ TTL / 3 max retries instead. RabbitMQ requires the consumer's
queue-declare arguments to be identical to the publisher's, and a queue
only receives what's published to the exchange it's bound to — so
`profile-worker` matches the **publisher's real behavior** (see
`cmd/profile-worker/main.go`). Orchestrator action needed: either rename
the exchange in api-service to `profile-tasks` (and fix its TTL) or update
CONTRACTS.md / ROUTING_KEYS.md to document `tasks-exchange` / 24h.

Exchanges are `direct` and durable. Every worker declares its own topology
idempotently on startup (main exchange, `<exchange>.dlx`, main queue with
`x-dead-letter-exchange`/`x-dead-letter-routing-key`/`x-message-ttl`, and
`<queue>.dlq`), and reconnects with backoff on connection/channel loss.
Messages that fail to parse or fail processing are nack'd **without**
requeue (so they route to the DLQ) instead of being redelivered forever.

Envelope: `{id, type, timestamp, payload, metadata}` — `metadata` is
optional and unknown top-level/payload fields are tolerated (forward
compatibility), per CONTRACTS.md §2.

## Environment variables

| Variable | Default | Notes |
|---|---|---|
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672/` | Primary connection string (CONTRACTS.md §4). |
| `RABBITMQ_HOST` / `RABBITMQ_PORT` / `RABBITMQ_USER` / `RABBITMQ_PASSWORD` / `RABBITMQ_VHOST` | `rabbitmq` / `5672` / `guest` / `guest` / _(default vhost)_ | Used only to build the connection string when `RABBITMQ_URL` is unset. |
| `HEALTH_PORT` | `8080` | Port for `/health` and `/ready`. |

## Manual smoke test

```bash
pip install aio-pika
RABBITMQ_HOST=localhost python tests/publish_test_messages.py
```

Publishes one MESSAGE_FORMAT.md-shaped message per queue. Start the
workers first so the exchanges/queues/DLQs exist and are bound — a
`direct` exchange silently drops a publish if nothing is bound to it.

## Docker

`Dockerfile.email` / `Dockerfile.image` / `Dockerfile.profile` are
multi-stage (`golang:1.24-alpine` builder → `alpine:3.19` runtime), run as
a non-root `worker` user, and are correct by inspection — not built in
this workflow (`docker build` is intentionally out of scope here).

## Known leftovers / risks

- Profile exchange contract drift — see table above.
- Very long simulated image work (`analyze` at `low` priority, ~31s) can
  outlive `BaseWorker.Shutdown`'s 10s HTTP-shutdown timeout budget; the
  consumer itself still waits unconditionally for the in-flight message to
  finish before closing the channel, so no message is lost, but
  orchestrators should set `terminationGracePeriodSeconds` generously
  (e.g. ≥35s) for image-worker.
- `internal/common/queue/publisher.go` (a `Publisher` type mirroring the
  consumer) is dead code — no worker publishes messages — and was left
  untouched; it has a pre-existing `NotifyPublish`-after-`Publish` ordering
  bug that only matters if something starts using it.
- No integration test against a live broker (none was available/started
  per task constraints); reconnect/DLQ topology was verified by matching
  the publisher's Go source line-for-line, not by observing RabbitMQ.
