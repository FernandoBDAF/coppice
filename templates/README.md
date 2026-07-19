# Templates (ADR-010.3 · phase v8)

The lab's original motivation, made deliverable: hardened pieces become
**copy-ready templates** a real project bootstraps from in under a day
(PRD success metric — EXP-80 measures it honestly). Rules:

- A template ships **only patterns experiments have beaten on** — every
  claim in a template README cites the experiment that proved it.
- Copy-then-own: no semver, no shared library. A piece graduates to its
  own repo on first real adoption — procedure in
  [GRADUATION.md](GRADUATION.md).
- Each template carries: `README.md` (adapt-in-a-day guide), a bootstrap
  test (`test/bootstrap.sh` — proves a fresh copy builds and passes its
  smoke), its k8s base + compose snippet, and its own minimal CI workflow
  runnable in the consuming repo.

| Template | Source of truth (extraction origin) | Status |
|---|---|---|
| [auth-service](auth-service/) | `auth-service/` post-v4 (JWKS, sessions, rotation, lockout, roles) | skeleton — v8-HANDOFF §1 |
| [worker-go](worker-go/) | `graph-worker/operational-workers/internal/common/*` post-v4 (envelope, consume loop, retry tiers, idempotency, results) | skeleton — v8-HANDOFF §2 |
| [api-publisher](api-publisher/) | `api-service` post-v4 (outbox+relay, typed submission, JWKS middleware) | skeleton — v8-HANDOFF §3 |
| [patterns/](patterns/) | docs, not code (cache-aside, storage pipeline, deploy shapes) | outlines — v8-HANDOFF §4 |

The one Helm chart (ADR-002.1) packages **worker-go**: `worker-go/chart/`
— the deliberate Helm-literacy exercise; kustomize remains the lab's tool.
