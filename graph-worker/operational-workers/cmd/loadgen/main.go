// loadgen (ADR-004.4) — AMQP load generator for queue-side experiments.
//
// SKELETON for the v4 handoff: flag surface and envelope construction are
// final; the publish loop is the remaining work (HANDOFF §B4). Keep it
// envelope-correct (graph-worker/shared/contracts/MESSAGE_FORMAT.md) so
// consumers treat generated load exactly like api-service traffic.
//
// Target CLI (wired as an experiment step type by the runner):
//
//	loadgen -url amqp://guest:guest@rabbitmq:5672/ \
//	        -routing-key email.send -rate 200 -duration 30s -confirm
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type opts struct {
	URL        string
	RoutingKey string
	Rate       int           // messages per second
	Duration   time.Duration // total run time
	Confirm    bool          // publisher confirms on/off (throughput vs safety demo)
	Payload    string        // optional JSON payload override
}

// envelope mirrors the publisher contract (id/type/timestamp/payload/metadata).
type envelope struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Timestamp string            `json:"timestamp"`
	Payload   json.RawMessage   `json:"payload"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// defaultPayloads produce a valid payload per routing key so workers
// process (not reject) generated load. Keep in sync with the processors'
// Validate rules.
var defaultPayloads = map[string]string{
	"email.send":     `{"email_type":"notification","recipient":"loadgen@lab.local","subject":"loadgen"}`,
	"image.process":  `{"operation":"resize","source_url":"minio://documents-raw/loadgen.png","target_path":"/tmp/out.png","width":64,"height":64}`,
	"profile.task":   `{"task_type":"sync","profile_id":"00000000-0000-0000-0000-000000000001"}`,
	"document.process": `{"document_id":"00000000-0000-0000-0000-000000000002","storage_path":"loadgen/x.txt","storage_bucket":"documents-raw"}`,
}

func main() {
	var o opts
	flag.StringVar(&o.URL, "url", "amqp://guest:guest@rabbitmq:5672/", "AMQP URL")
	flag.StringVar(&o.RoutingKey, "routing-key", "email.send", "routing key (also envelope type)")
	flag.IntVar(&o.Rate, "rate", 100, "messages per second")
	flag.DurationVar(&o.Duration, "duration", 30*time.Second, "run duration")
	flag.BoolVar(&o.Confirm, "confirm", true, "wait for publisher confirms")
	flag.StringVar(&o.Payload, "payload", "", "JSON payload override (default: per-routing-key valid payload)")
	flag.Parse()

	if _, ok := defaultPayloads[o.RoutingKey]; !ok && o.Payload == "" {
		fmt.Fprintf(os.Stderr, "unknown routing key %q and no -payload given\n", o.RoutingKey)
		os.Exit(2)
	}

	// TODO(v4, HANDOFF §B4): implement run(o):
	//  - dial amqp091, open channel, Confirm(false) when o.Confirm
	//  - resolve exchange from routing key via the definitions-derived map
	//    (profile-tasks/email-tasks/image-tasks/document-tasks — do NOT
	//    declare topology; broker owns it, ADR-008.4)
	//  - ticker at o.Rate msg/s for o.Duration; per message: fresh uuid id,
	//    RFC3339 UTC timestamp, metadata{source:"loadgen"}; persistent JSON
	//  - stdout summary: sent, confirmed, nacked, errors, effective rate
	//  - exit 1 on any nack/unroutable when -confirm (mandatory publish)
	fmt.Fprintln(os.Stderr, "loadgen: publish loop not implemented yet (v4 handoff §B4)")
	os.Exit(3)
}
