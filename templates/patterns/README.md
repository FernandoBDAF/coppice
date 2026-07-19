# Pattern docs (not code) — OUTLINES for v8

Each becomes a standalone doc when v8 executes (v8-HANDOFF §4). The
outline pins scope + the experiment evidence each must cite.

## cache-aside.md
Postgres-backed cache-aside with Redis, as api-service does it. MUST
cover: the EXP-12 finding verbatim (Redis outage → 0% errors but silent
0.3–4s dial penalties; only /ready reports it) + the mitigation menu
(dial timeout tuning, fail-fast flag, per-call circuit); stampede caveats
(no singleflight in the lab — say so and show the seam); key naming +
invalidation discipline (`profile:<id>`, list keys, delete-on-write).

## storage-pipeline.md
Presigned-URL object pipeline (upload → object store → metadata row →
async processing → status lifecycle). MUST cover: compensating delete on
DB failure, publish-after-commit via outbox (EXP-42), status write-back
loop (task-results), presigned expiry choice (15m) and S3-vs-MinIO
presign behavior diffs (v5 EXP-50 evidence), orphan reconciliation sweep
design (objects without rows / rows without objects).

## deploy-shapes.md
The kustomize layout that survived v2→v5: base + overlays, single-source
configMapGenerator trick (LoadRestrictionsNone), probes/limits/PDB
defaults per service class, netpol starter (default-deny + dns +
explicit allows; union-of-allows lesson from EXP-23), drift-check as a
CI institution.
