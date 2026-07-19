# Template: auth-service — SKELETON

**Extract from** (post-v4 `auth-service/`): JWKS (RS256, kid, rotation),
sessions table + one-time refresh rotation + reuse detection, account
lockout, roles middleware, metrics/health, seed/init scripts. **Trim:**
lab-specific rate-limit envs (EXP-08 apparatus), OpenAPI docs generation
(keep the zod schemas), audit-log niceties beyond the auth events.

**Ships as:** the trimmed service + migrations, config surface doc (every
env var: name, default, constraint), k8s base + compose snippet, bootstrap
test, CI workflow (typecheck+build+test).

**Proven by:** EXP-02 (contract), EXP-08→EXP-43 (the introspection SPOF
and its JWKS retirement — the template defaults to JWKS-only and documents
introspection as an opt-in), session-reuse tests (v4 A7).

## How to adapt
1. Copy; rename service + DB; run migrations; generate keys (init script
   included).
2. Claims: edit ONE place (`TokenService.basePayload`) — document any
   claim your API consumers verify.
3. `test/bootstrap.sh`: register→login→validate→refresh-rotate→reuse
   rejected→logout, against compose postgres. Green = adapted.

Details + step list: v8-HANDOFF §1.
