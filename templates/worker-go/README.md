# Template: async worker (Go) — SKELETON

**Extract from** (post-v4 `graph-worker/operational-workers/`):
`internal/common/queue/` (envelope, consume loop w/ reconnect+backoff,
retry-tier publishing, passive topology verification),
`internal/common/idempotency/` (SETNX guard), `internal/common/base/`
(worker lifecycle, health server, metrics), one example processor
(email, trimmed), results publishing (task-results), plus a
`definitions.json` **fragment generator** for a new queue (a `--pipeline
name:rk` flag on `scripts/rabbitmq/generate-definitions.py`, copied in).

**Ships as:** `worker-go/{cmd/example-worker,internal/...}` with module
path `example.com/worker` (the adapt step renames it), Dockerfile, compose
snippet, k8s base, `chart/` (the ADR-002.1 Helm exercise), bootstrap test.

**Proven by:** EXP-04/05/06 (durability + DLQ), EXP-40 (retry tiers),
EXP-41 (idempotency), EXP-42 (results loop) — README must cite these with
one-line summaries when the template is cut.

## How to adapt (target: < 1 hour of the < 1 day budget)
1. Copy dir; rename module (`go mod edit -module`), worker type, queue
   names; run the fragment generator for your queue → merge into your
   broker definitions.
2. Replace the example processor (implement `MessageProcessor`).
3. `test/bootstrap.sh` green (build + unit tests + compose smoke:
   publish 10, consume 10, poison 1 → DLQ 1).

## Bootstrap test contract (v8-HANDOFF §5)
`test/bootstrap.sh`: docker compose up (rabbitmq+redis+worker) → seed →
assert counts via management API → down. Deliberately breakable (EXP-81:
a broken copy must fail it).
