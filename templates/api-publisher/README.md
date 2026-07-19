# Template: API publisher (Go) — SKELETON

**Extract from** (post-v4 `api-service/`): outbox table + relay
(`internal/infrastructure/outbox/`), typed task submission
(`internal/domain/task/` envelope construction + rk→exchange map), JWKS
verification middleware (`internal/infrastructure/auth/jwks_verifier.go`),
plus a `CONTRACTS.md` skeleton (topology table, envelope, auth contract —
the doc pattern that kept this lab honest).

**Trim:** documents/profiles domains, MinIO, Redis cache — this template
is "an API that publishes work reliably", nothing more.

**Proven by:** EXP-42 (outbox crash-consistency), EXP-43 (local verify
under auth outage), EXP-45 (scored smoke as the consuming repo's CI seed).

## How to adapt
1. Copy; rename module; define your task types in one file (typed payload
   + rk); point the fragment generator at your queues.
2. Wire your domain writes through `outbox.Add` in the same tx.
3. `test/bootstrap.sh`: POST a task → outbox row → relay publishes →
   management API shows the message; kill-between-commit-and-publish
   variant proves the crash window is closed.

Details: v8-HANDOFF §3.
