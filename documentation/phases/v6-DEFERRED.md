# Phase v6 — deferred validation & follow-ups

**Context:** v6 was implemented in an expedited pass (2026-07-19): the
`lab-controld` control plane (registry loader with SIGHUP reload, action
execution via `sh -c` from the repo root, SSE streaming, JSONL run history,
the one-action-per-(system,target) 409 guard, the destructive-verb confirm
gate), the ADR-005.4 auth gate, the wave-2 endpoints (experiments catalog,
outcome recording, session recorder), and the Mission Control cockpit UI all
landed and were unit-tested with fixture registries and fake (`sh -c 'echo
ok'` / `'exit 2'`) commands. But the exit experiments (EXP-60..63) need a
**running compose/kind stack and a browser**, and none ran this pass. This
file is the honest ledger of what remains before `lab-v6.0` can be tagged.

## Must run before tagging lab-v6.0

| Item | What to do | Acceptance check |
|---|---|---|
| EXP-60 terminal-free session | Fresh browser → kind target → launch lab → EXP-04 guided (panels inline) → induce a v4 chaos fault + diagnose → record the outcome, zero terminal | `GET /api/sessions/{id}/summary` renders a paste-ready write-up; it lands in `documentation/experiments/`; the runs log shows every step delegated to make (EXPERIMENTS.md §EXP-60) |
| EXP-61 target parity | Same lab card up/down on compose and kind; aws chip shows correct state with a session up; force a failing make and confirm the action fails loudly | Actions stream real stdout; a non-zero make exit → `state:failed` with the exit code, never a silent success (EXPERIMENTS.md §EXP-61) |
| EXP-62 library round-trip | Run a scored experiment from the UI, then `make experiment E=<id>` for the same id | Same pass/fail both ways; outcome appended under `documentation/experiments/` (EXPERIMENTS.md §EXP-62) |
| EXP-63 control-plane safety | Localhost-bound daemon refuses a remote connection; with `CONTROLD_ENABLE_AWS=1` + `CONTROLD_TOKEN`, hit `/api/*` with a wrong token | Remote connect refused on the localhost bind; wrong token → 401 **and** an audit log line (EXPERIMENTS.md §EXP-63) |
| "Every action delegates to make" spot-check | Drive a handful of UI actions, read `mission-control/controld/runs/YYYY-MM-DD.jsonl` | Every `ActionRecord.command` is a `make …` (or registry) invocation — no hidden shell (phase-doc acceptance) |
| hello-guest launchable from the UI | Launch hello-guest from its system card via its `systems/hello-guest.yaml` entry | Guest comes up on the target; proves the v7-readiness of the systems model (phase-doc acceptance) |

## Seams on parallel work — reconciled 2026-07-19

The v5 execution (which carries v4-final) was merged into this branch
(`Merge phase/v5`), closing the two seams the first pass had to stub:

- **Scored runs now attach real per-assertion results.** The real runner
  (`scripts/experiments/run.py`, ADR-004.1) is on this branch with all 12
  `experiments/*.yaml`. For `verb:experiment` actions controld sets
  `EXPERIMENT_REPORT_DIR` to a per-action dir under `runs/reports/` and, on
  completion, parses the runner's junit XML into `ActionRecord.report`
  (`{passed,total,failed,assertions[]}`); the UI renders the breakdown in the
  action modal, library result, and a compact history badge. Pass/fail stays
  exit-code driven; a missing/unparseable report degrades to no breakdown,
  never a fake result. **Remaining check:** EXP-62's UI result == CLI result
  for a passing id, on a live stack.
- **AWS availability is a live probe.** `AWSTargetEntry(cfg)` runs read-only
  `terraform -chdir=deploy/aws/session output -raw cluster_name` (10s timeout,
  60s cache) — session up → the aws chip enables with the cluster name; any
  failure (no terraform, no init, no state) → disabled with the honest reason.
  Live aws parity (EXP-61's aws leg) and a token+TLS drill against a real
  session still need a real session up.
- **The PR #15 graphrag hotfix is on `main` but not yet in this stack**
  (it landed after `phase/v5`'s v4-final merge). EXP-01's cold start can
  crash-loop graphrag without it; reconcile `main` bottom-up (or rebase the
  stack) before scoring EXP-01 from the UI.
- **hello-guest cards read "unknown".** The read API tracks only the lab; a
  guest card shows state `unknown` until you run its `status` verb as an
  on-demand action ("Check status"). **Check:** the status action returns the
  guest's real state and the card updates.
- **No EXP-60..63 live run this pass.** All four need a running compose/kind
  stack and a browser; the code was exercised only against unit fixtures.
  **Check:** the table above, each row green with a write-up.

## Known caveats shipped knowingly

- **`runs/` is gitignored, per-day JSONL, no database** (ADR-005.2 — no new
  storage engine). History lives at
  `mission-control/controld/runs/YYYY-MM-DD.jsonl`; losing it loses history,
  not state — the lab's truth is always the live target.
- **Auth is off on localhost, by decision (ADR-005.4).** The default
  127.0.0.1 bind stays no-auth for zero-friction local use. Bearer auth
  engages only when `CONTROLD_TOKEN` is set; enabling the aws target
  (`CONTROLD_ENABLE_AWS=1`) *requires* a token **and** TLS
  (`CONTROLD_TLS_CERT`/`CONTROLD_TLS_KEY`) or the daemon refuses to boot —
  a hard gate before any remote reach, not a runtime warning.
- **The registry is the whole action whitelist.** Nothing outside
  `systems/*.yaml` is invokable; `{n}` is validated 1..10 and the experiment
  id against `^exp-[a-z0-9-]+$` before any placeholder reaches the shell
  (v6-HANDOFF §2). Destructive verbs (down, resets) require
  `params.confirm="true"`.
- **SSE, not WebSocket, for streaming** (actions.go): one-way stdout/stderr
  fanout, `EventSource` is enough; the SSE query-string carries `?token=` when
  auth is on (browsers can't set a Bearer header on `EventSource`).

## Nice-to-haves consciously skipped

- Embedded Grafana panels beyond iframe/deep links from each system's `links`
  map (charts stay in Grafana — v6-HANDOFF §4).
- Any UI framework beyond Next/React; no client-side charting library.
- Multi-user / RBAC anything (ADR-005.4 — single-operator local tool);
  loam integration stays patterns-only (ADR-005.5 — loam onboards as a v7
  guest system, not UI plumbing).

## Tag policy

`lab-v6.0` is tagged **only after EXP-60..63 pass and are written up** (EXP-60's
session record is the showcase). The acceptance checkboxes in
[v6-mission-control.md](v6-mission-control.md) stay unchecked until then —
they gate the tag.
