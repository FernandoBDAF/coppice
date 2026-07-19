#!/usr/bin/env python3
"""Scored experiment runner (ADR-004.1) — SKELETON for the v4 handoff.

Contract (final; implementation pending — HANDOFF §B2):
  make experiment E=exp-04   →  python3 scripts/experiments/run.py exp-04
  make experiments           →  python3 scripts/experiments/run.py --list

Behavior to implement:
 1. Load experiments/<id>.yaml (schema: experiments/README.md). Reject empty
    `assertions` (a scored run must be falsifiable — EXP-45).
 2. Pre-check `needs` (compose project up / kind cluster present / obs stack
    / guest running) and fail fast with a pointed message.
 3. Run `steps` sequentially via subprocess (shell=True, fail-fast,
    background: true → Popen without wait; track for cleanup).
 4. Poll `assertions` every 5s until each passes or its timeout lapses:
      promql — GET {PROM_URL}/api/v1/query, compare first sample vs op/value
      http   — status equality (optionally jq-path equality later)
      cli    — exit code 0
 5. Always run `cleanup`. Exit 0 iff all assertions passed.
 6. Append a result block (id, date, pass/fail per assertion, durations) to
    documentation/experiments/RUNS.md and emit junit-ish XML to
    .experiment-results/<id>-<timestamp>.xml (gitignored) for CI phase 2.

Keep dependencies to stdlib + PyYAML (already a CI dep for drift-check).
"""
import sys


def main() -> int:
    if len(sys.argv) < 2:
        print(__doc__)
        return 2
    print(
        f"runner skeleton: '{sys.argv[1]}' not executed — implement per "
        "documentation/phases/v4-HANDOFF.md §B2 (schema: experiments/README.md)",
        file=sys.stderr,
    )
    return 3


if __name__ == "__main__":
    sys.exit(main())
