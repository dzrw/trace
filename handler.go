package trace

import "time"

type Metric string
type HandlerCaps int16

const (
	SupportsLogs    HandlerCaps = 1 << iota
	SupportsMetrics HandlerCaps = 1 << iota
	SupportsTraces  HandlerCaps = 1 << iota
)

// A Handler processes traces and associated event logs.
type Handler interface {
	// Capabilities reports whether this handler supports logs, metrics,
	// and traces.
	Capabilities() HandlerCaps

	// Enabled returns whether this handler accepts logs at a level.
	Enabled(Level) bool

	// Log an event.
	Log(Level, [][]Attr) error

	// Count records a delta to a counter.
	Count(Metric, int64) error

	// Gauge records the value of a gauge.
	Gauge(Metric, int64) error

	// Duration records an elapsed time.
	Duration(Metric, time.Duration) error
}
