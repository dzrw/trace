package trace

import (
	"runtime"
	"time"

	"github.com/segmentio/ksuid"
)

type Tracer interface {
	ID() ksuid.KSUID // ID returns the unique identifier for the trace
	Name() string    // Name returns the title of the trace
	Finish()         // Finish closes the trace
}

type Trace[T ProbeC] interface {
	Tracer

	// Error writes an event log associated with a span.
	Error(probe T, err error, attrs ...Attr)

	// Warn writes an event log associated with a span.
	Warn(probe T, attrs ...Attr)

	// Info writes an event log associated with a span.
	Info(probe T, attrs ...Attr)

	// Debug writes an event log associated with a span.
	Debug(probe T, attrs ...Attr)

	// Log an event associated with a span. Log events are immediately flushed
	// to the Logger associated with the package. If telemetry is not configured,
	// this operation is a nop.
	Log(probe T, level Level, attrs ...Attr)

	// Count updates a counter associated with a span. Counters are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Count(probe T, delta int64) int64

	// Gauge updates a gauge associated with a span. Gauges are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Gauge(probe T, value int64)

	// Histogram updates a log-linear histogram associated with a span. Histograms are
	// flushed when the span is closed. If telemetry is not configured, this operation
	// is a nop.
	Histogram(probe T, sample int64)
}

// New creates and returns a Span.
func New[T ProbeC](pkg Package, name string) Trace[T] {
	tr := &traceimpl[T]{
		pkg:  pkg,
		id:   ksuid.New(),
		name: name,
		then: time.Now(),
	}

	if pkg.CaptureSourceInfo() {
		if pc, file, line, ok := runtime.Caller(1); ok {
			tr.file = file
			tr.line = line
			if f := runtime.FuncForPC(pc); f != nil {
				tr.funcName = f.Name()
			}
		}

		tr.gid = __caution__GetGoroutineID()
	}

	return tr
}

type traceimpl[T ProbeC] struct {
	pkg      Package
	id       ksuid.KSUID
	name     string
	then     time.Time
	duration time.Duration

	gid      uint64 // goroutine identifier
	file     string
	line     int
	funcName string

	counters map[Probe]int64
	gauges   map[Probe]int64
}

func (tr *traceimpl[T]) ID() ksuid.KSUID {
	return tr.id
}

func (tr *traceimpl[T]) Name() string {
	return tr.name
}

func (tr *traceimpl[T]) Finish() {
	tr.duration = time.Since(tr.then)
	if h := tr.pkg.Handler(); h != nil {
		for p, val := range tr.counters {
			h.Count(tr, p, val)
		}
		for p, val := range tr.gauges {
			h.Gauge(tr, p, val)
		}
	}
}

func (tr *traceimpl[T]) Error(probe T, err error, attrs ...Attr) {
	attrs = append(attrs, Any("error", err))
	tr.Log(probe, ErrorLevel, attrs...)
}

func (tr *traceimpl[T]) Warn(probe T, attrs ...Attr) {
	tr.Log(probe, WarnLevel, attrs...)
}

func (tr *traceimpl[T]) Info(probe T, attrs ...Attr) {
	tr.Log(probe, InfoLevel, attrs...)
}

func (tr *traceimpl[T]) Debug(probe T, attrs ...Attr) {
	tr.Log(probe, DebugLevel, attrs...)
}

func (tr *traceimpl[T]) Log(probe T, level Level, attrs ...Attr) {
	if level <= probe.Level() {
		if h := tr.pkg.Handler(); h != nil && h.Enabled(level) {
			var file = ""
			var line = 0
			h.Log(tr, NewEventLog(time.Now(), level, probe.String(), file, line))
		}
	}
}

func (tr *traceimpl[T]) Count(probe T, delta int64) (val int64) {
	val = delta
	if v, ok := tr.counters[probe]; ok {
		val = val + v
	}
	tr.counters[probe] = val
	return
}

func (tr *traceimpl[T]) Gauge(probe T, value int64) {
	tr.gauges[probe] = value
}

func (tr *traceimpl[T]) Histogram(probe T, sample int64) {
	// TODO: https://github.com/openhistogram/libcircllhist
}
