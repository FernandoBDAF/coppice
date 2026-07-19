package queue

// Retry-tier routing (ADR-008.1) — SKELETON for the v4 handoff.
//
// The broker owns the wait-queues (deploy/rabbitmq/definitions.json):
//   <queue>.retry.{5s,30s,2m}, bound to <exchange>.retry on
//   <rk>.retry.<tier>, each dead-lettering back to the main exchange.
//
// Consumer contract (replaces nack-to-DLQ for retryable errors):
//  1. attempt = DeathCount(headers) — how many times this message already
//     cycled through a wait-queue (x-death header, count for the main queue).
//  2. tier, ok := NextRetryTier(attempt); if ok → publish the message body
//     unchanged to "<exchange>.retry" with routing key rk+".retry."+tier,
//     persistent, then ACK the delivery.
//  3. if !ok (tiers exhausted) → publish to "<exchange>.dlx" with rk (poison
//     path), then ACK. The DLQ is reached explicitly, never via broker nack.
// Unretryable errors (unmarshal/validation) skip straight to step 3.
//
// documentation/phases/v4-HANDOFF.md §A3 has the consumer.go integration
// diff and the x-death parsing details (amqp.Table nesting).

// RetryTiers in escalation order; must match the generator
// (scripts/rabbitmq/generate-definitions.py RETRY_TIERS).
var RetryTiers = []string{"5s", "30s", "2m"}

// NextRetryTier maps how-many-times-already-retried to the next wait-queue
// suffix. attempt 0 → "5s", 1 → "30s", 2 → "2m", ≥3 → exhausted.
func NextRetryTier(attempt int) (string, bool) {
	if attempt < 0 || attempt >= len(RetryTiers) {
		return "", false
	}
	return RetryTiers[attempt], true
}

// DeathCount extracts the number of prior retry cycles from the x-death
// header. TODO(v4): implement against amqp.Table — x-death is a []interface{}
// of tables; sum the "count" entries whose "queue" is one of our
// <queue>.retry.* wait-queues (counting the work queue itself would
// double-count expiry loops). Unit-test with a captured header fixture.
func DeathCount(headers map[string]any) int {
	_ = headers // TODO(v4) — see HANDOFF §A3
	return 0
}
