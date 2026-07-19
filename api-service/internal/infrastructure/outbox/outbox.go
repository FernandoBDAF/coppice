// Package outbox implements the transactional outbox (ADR-008.3):
// domain writes and their task events commit in ONE api_db transaction;
// a relay goroutine publishes pending rows with confirms and marks them
// sent. This closes the EXP-42 crash window (DB committed, publish lost).
//
// SKELETON for the v4 handoff — types and SQL are final, method bodies and
// the migration are the remaining work (documentation/phases/v4-HANDOFF.md
// §A5, which also contains the exact CREATE TABLE migration to add under
// api-service/migrations/).
package outbox

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNotImplemented = errors.New("outbox: not implemented (v4 handoff)")

// Row is one to-be-published task event. Envelope fields are stored
// pre-serialized (payload JSON + routing key) so the relay never needs
// domain knowledge.
type Row struct {
	ID         int64      `db:"id"`
	RoutingKey string     `db:"routing_key"`
	Envelope   []byte     `db:"envelope"` // full task.Message JSON, built at write time
	CreatedAt  time.Time  `db:"created_at"`
	SentAt     *time.Time `db:"sent_at"`
	Attempts   int        `db:"attempts"`
}

// Store persists events inside the caller's transaction — the whole point:
// callers pass the SAME tx that writes the domain row.
type Store struct{ db *sqlx.DB }

func NewStore(db *sqlx.DB) *Store { return &Store{db: db} }

// Add inserts a pending row within tx.
// SQL: INSERT INTO outbox (routing_key, envelope) VALUES ($1, $2)
func (s *Store) Add(ctx context.Context, tx *sqlx.Tx, routingKey string, envelope []byte) error {
	_ = ctx
	_ = tx
	_ = routingKey
	_ = envelope
	return ErrNotImplemented // TODO(v4) HANDOFF §A5
}

// PendingBatch fetches unsent rows FOR UPDATE SKIP LOCKED (relay is safe to
// run on every replica).
func (s *Store) PendingBatch(ctx context.Context, limit int) ([]Row, error) {
	_ = ctx
	_ = limit
	return nil, ErrNotImplemented // TODO(v4) HANDOFF §A5
}

// MarkSent stamps sent_at; MarkFailed bumps attempts (relay backs off and
// alerts on attempts > threshold — rows are never deleted by the relay).
func (s *Store) MarkSent(ctx context.Context, ids []int64) error {
	_ = ctx
	_ = ids
	return ErrNotImplemented // TODO(v4) HANDOFF §A5
}

func (s *Store) MarkFailed(ctx context.Context, ids []int64) error {
	_ = ctx
	_ = ids
	return ErrNotImplemented // TODO(v4) HANDOFF §A5
}

// Publisher is satisfied by the existing rabbitmq publisher (confirm-mode).
type Publisher interface {
	PublishRaw(ctx context.Context, routingKey string, body []byte) error
}

// Relay drains the outbox: poll (or LISTEN/NOTIFY later) → PendingBatch →
// publish with confirms → MarkSent/MarkFailed. Runs from cmd/server/main.go
// as a goroutine; interval ~250ms idle, immediate loop while batch full.
// TODO(v4) HANDOFF §A5: implement loop + prometheus metrics
// (api_outbox_pending gauge, api_outbox_published_total,
// api_outbox_publish_errors_total) + graceful shutdown.
func Relay(ctx context.Context, store *Store, pub Publisher, interval time.Duration) error {
	_ = store
	_ = pub
	_ = interval
	<-ctx.Done()
	return ErrNotImplemented
}
