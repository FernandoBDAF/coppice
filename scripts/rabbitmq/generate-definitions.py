#!/usr/bin/env python3
"""Generate the broker-owned topology (ADR-008.4).

Single source of truth for exchanges/queues/bindings. Emits:
  deploy/rabbitmq/definitions.json        — loaded by the broker at boot
  deploy/rabbitmq/ROUTING_KEYS.generated.md — human contract table

Design (ADR-008.1/.3/.5/.6):
- Work queues have NO TTL (email is the exception: staleness TTL routes to
  the distinct `email-expired` queue via routing key `email.expired`, never
  the poison DLQ).
- Consumers do not rely on broker retry args: on retryable failure they
  PUBLISH to `<exchange>.retry` with routing key `<rk>.retry.<tier>` and ack;
  wait-queues dead-letter back to the main exchange after their TTL.
  Tiers: 5s → 30s → 2m, chosen by x-death count; after the last tier the
  consumer publishes to `<exchange>.dlx` with `<rk>` (poison path) and acks.
- The queue-level DLX therefore only fires for TTL expiry (email) and for
  legacy nacks (unmarshal poison), both of which are still safe.
- `task-results` closes the loop: workers/graphrag publish completion or
  failure; api-service consumes and advances document/task status.
- No default-tasks parking lot (ADR-008.6).

Run: python3 scripts/rabbitmq/generate-definitions.py
"""
import base64
import hashlib
import json
import os

# (rk, exchange, queue, dlq_ttl_ms, prefetch, consumer)
PIPELINES = [
    ("profile.task", "profile-tasks", "profile-processing", 86_400_000, 2, "profile-worker"),
    ("email.send", "email-tasks", "email-processing", 86_400_000, 5, "email-worker"),
    ("image.process", "image-tasks", "image-processing", 259_200_000, 1, "image-worker"),
    ("document.process", "document-tasks", "document-processing", 604_800_000, 1, "graphrag-service"),
]
RETRY_TIERS = [("5s", 5_000), ("30s", 30_000), ("2m", 120_000)]
EMAIL_STALENESS_TTL_MS = 3_600_000  # 1h — same envelope as the old queue TTL
EMAIL_EXPIRED_RETENTION_MS = 86_400_000

ROOT = os.path.join(os.path.dirname(__file__), "..", "..")
OUT_DIR = os.path.join(ROOT, "deploy", "rabbitmq")


def guest_password_hash() -> str:
    # rabbit_password_hashing_sha256: base64(salt + sha256(salt + password)).
    # Fixed salt so the file is reproducible; lab-only credentials (guest).
    salt = bytes.fromhex("cafebabe")
    return base64.b64encode(salt + hashlib.sha256(salt + b"guest").digest()).decode()


def build():
    exchanges, queues, bindings = [], [], []

    def exchange(name):
        exchanges.append({"name": name, "vhost": "/", "type": "direct",
                          "durable": True, "auto_delete": False, "internal": False,
                          "arguments": {}})

    def queue(name, args):
        queues.append({"name": name, "vhost": "/", "durable": True,
                       "auto_delete": False, "arguments": args})

    def bind(src, dest, rk):
        bindings.append({"source": src, "vhost": "/", "destination": dest,
                         "destination_type": "queue", "routing_key": rk,
                         "arguments": {}})

    for rk, ex, q, dlq_ttl, _prefetch, _consumer in PIPELINES:
        exchange(ex)
        exchange(f"{ex}.dlx")
        exchange(f"{ex}.retry")

        work_args = {
            "x-dead-letter-exchange": f"{ex}.dlx",
            # expiry-vs-poison split (ADR-008.5): consumers publish poison
            # explicitly, so the queue's own dead-letter path is expiry-only
            "x-dead-letter-routing-key": "email.expired" if rk == "email.send" else rk,
        }
        if rk == "email.send":
            work_args["x-message-ttl"] = EMAIL_STALENESS_TTL_MS
        queue(q, work_args)
        bind(ex, q, rk)

        for tier, ttl in RETRY_TIERS:
            rq = f"{q}.retry.{tier}"
            queue(rq, {
                "x-message-ttl": ttl,
                "x-dead-letter-exchange": ex,
                "x-dead-letter-routing-key": rk,
            })
            bind(f"{ex}.retry", rq, f"{rk}.retry.{tier}")

        queue(f"{q}.dlq", {"x-message-ttl": dlq_ttl})
        bind(f"{ex}.dlx", f"{q}.dlq", rk)

    # email staleness parking (never mixed with poison)
    queue("email-expired", {"x-message-ttl": EMAIL_EXPIRED_RETENTION_MS})
    bind("email-tasks.dlx", "email-expired", "email.expired")

    # completion/failure feedback loop (ADR-008.3)
    exchange("task-results")
    queue("task-results", {})
    bind("task-results", "task-results", "task.result")

    return {
        "rabbit_version": "3.12.0",
        "vhosts": [{"name": "/"}],
        # load_definitions at boot skips default-user creation, so guest must
        # be declared here (lab-only credential, matches compose/cluster)
        "users": [{"name": "guest", "password_hash": guest_password_hash(),
                   "hashing_algorithm": "rabbit_password_hashing_sha256",
                   "tags": ["administrator"]}],
        "permissions": [{"user": "guest", "vhost": "/",
                         "configure": ".*", "write": ".*", "read": ".*"}],
        "topic_permissions": [],
        "parameters": [],
        "global_parameters": [],
        "policies": [],
        "exchanges": exchanges,
        "queues": queues,
        "bindings": bindings,
    }


def routing_table_md(defs):
    lines = [
        "# Routing keys & topology (GENERATED — do not edit)",
        "",
        "Source of truth: `scripts/rabbitmq/generate-definitions.py` →",
        "`deploy/rabbitmq/definitions.json` (ADR-008.4). Services verify",
        "(passively declare-check) but never author topology.",
        "",
        "| routing key | exchange | work queue | retry tiers | DLQ (TTL) | prefetch | consumer |",
        "|---|---|---|---|---|---|---|",
    ]
    for rk, ex, q, dlq_ttl, prefetch, consumer in PIPELINES:
        tiers = "/".join(t for t, _ in RETRY_TIERS)
        lines.append(f"| `{rk}` | `{ex}` | `{q}` | {tiers} | `{q}.dlq` ({dlq_ttl // 3_600_000}h) | {prefetch} | {consumer} |")
    lines += [
        "",
        f"- `email-processing` staleness TTL {EMAIL_STALENESS_TTL_MS // 60_000}min → `email-expired` (rk `email.expired`), never the poison DLQ (ADR-008.5).",
        "- Retry flow (ADR-008.1): consumer publishes to `<exchange>.retry` rk `<rk>.retry.<tier>`, acks; wait-queue TTL dead-letters back to the main exchange. After the last tier, consumer publishes to `<exchange>.dlx` rk `<rk>` (poison) and acks.",
        "- `task-results` (rk `task.result`): workers/graphrag publish completion/failure; api-service consumes (ADR-008.3).",
        "- No `default-tasks` fallback — unknown routing keys are a publisher bug and fail fast (ADR-008.6).",
        "",
        f"Totals: {len(defs['exchanges'])} exchanges, {len(defs['queues'])} queues, {len(defs['bindings'])} bindings.",
    ]
    return "\n".join(lines) + "\n"


def main():
    os.makedirs(OUT_DIR, exist_ok=True)
    defs = build()
    with open(os.path.join(OUT_DIR, "definitions.json"), "w") as f:
        json.dump(defs, f, indent=2, sort_keys=False)
        f.write("\n")
    with open(os.path.join(OUT_DIR, "ROUTING_KEYS.generated.md"), "w") as f:
        f.write(routing_table_md(defs))
    print(f"wrote {len(defs['exchanges'])} exchanges, {len(defs['queues'])} queues, "
          f"{len(defs['bindings'])} bindings")


if __name__ == "__main__":
    main()
