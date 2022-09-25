package trace

import "time"

type HandlerFlags int16

const (
	FlagSourceInfo  HandlerFlags = 1 << iota
	FlagGoroutineID HandlerFlags = 1 << iota
)

// A Handler processes traces and associated event logs.
type Handler interface {
	// Flags returns the options set on this handler.
	Flags() HandlerFlags

	// Enabled returns whether this handler accepts logs at a level.
	Enabled(Level) bool

	// Count records a delta to a counter.
	Count(Tracepoint, int64) error

	// Gauge records the value of a gauge.
	Gauge(Tracepoint, int64) error

	// Duration records an elapsed time.
	Duration(Tracepoint, time.Duration) error

	// Histogram records a sample.
	Histogram(Tracepoint, int64) error

	// Log an event.
	Log(Trace, Level, ...[]Attr) error

	TraceCreated(Trace, []Attr)
	TraceFinished(Trace, []Attr)
}
