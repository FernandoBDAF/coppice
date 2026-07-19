# Scored experiment definitions (ADR-004.2 · phase v4)

One YAML per experiment; EXPERIMENTS.md stays the prose companion (Goal /
Watch / narrative) and becomes a generated index in v4. The runner
(`make experiment E=<id>` → `scripts/experiments/run.py`) executes `steps`,
then polls `assertions` until all pass or time out — exit 0/1 with a
junit-ish report appended to `documentation/experiments/`.

## Schema v0

```yaml
id: exp-04                # matches EXPERIMENTS.md heading
title: Burst absorption & drain
needs: [compose]          # compose | kind | obs | guest:<name> — runner pre-checks
steps:                    # executed sequentially, shell, fail-fast
  - run: make queues
  - run: make sim-burst
    background: false     # true → fire and continue (loadgen, floods)
watch:                    # not executed — prose refs the runner prints
  - "Lab Overview → Queue depth"
assertions:               # polled (interval 5s) until pass or timeout
  - type: promql          # instant query, compare against threshold
    query: sum(rabbitmq_queue_messages{queue=~".*-processing"})
    op: "<="              # == != < <= > >=
    value: 0
    timeout: 300s         # keep polling until this deadline
  - type: http            # status (+ optional jq-path equality)
    url: http://localhost:8080/ready
    status: 200
    timeout: 30s
  - type: cli             # exit code 0 == pass; retried like the others
    run: docker compose ps --status running --services | grep -q email-worker
    timeout: 30s
cleanup:                  # always runs, even on failure
  - run: "true"
```

Conventions: every assertion needs `timeout`; `promql` needs the stack's
Prometheus (`PROM_URL`, default http://localhost:9090); ids are kebab-case
`exp-NN`. The runner treats an empty `assertions` list as a config error —
a scored experiment must be falsifiable (EXP-45 proves the runner rejects
rubber stamps).

## Status

- [x] schema drafted (this file) · [x] first migration: `exp-02.yaml`
- [ ] runner implementation — `scripts/experiments/run.py` (HANDOFF §B2)
- [ ] migrate EXP-01..12 (HANDOFF §B1 has the per-experiment notes)
- [ ] EXPERIMENTS.md generated index (HANDOFF §B1)
