package queue

import "time"

type Config struct {
	// Connection Settings
	URL       string
	Heartbeat time.Duration
	Locale    string

	// Exchange and Queue Settings
	Exchange   string
	Queue      string
	RoutingKey string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool

	// Dead-letter topology. These MUST stay argument-identical to whatever
	// the publisher (api-service) declares for the same routing key, or
	// RabbitMQ raises a 406 PRECONDITION_FAILED channel error on redeclare.
	// DLX is always "<Exchange>.dlx" and DLQ is always "<Queue>.dlq".
	MessageTTL    time.Duration // x-message-ttl on the main queue
	DeadLetterTTL time.Duration // x-message-ttl on the .dlq queue
	MaxRetries    int           // x-max-retries (informational, mirrors publisher)

	// Consumer Settings
	PrefetchCount int
	PrefetchSize  int
	Global        bool

	// Publisher Settings
	Mandatory bool
	Immediate bool

	// Reconnect Settings
	RetryDelay time.Duration // initial backoff between reconnect attempts

	// Logging
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		Durable:       true,
		AutoDelete:    false,
		Exclusive:     false,
		NoWait:        false,
		PrefetchCount: DefaultPrefetchCount,
		MaxRetries:    DefaultMaxRetries,
		RetryDelay:    DefaultRetryDelay,
		Heartbeat:     DefaultHeartbeat,
		Locale:        DefaultLocale,
		LogLevel:      "info",
	}
}
