package trace

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Trace interface {
	ID() ksuid.KSUID        // ID returns the unique identifier for the trace
	Elapsed() time.Duration // Elapsed returns the time elapsed since the trace started.

	// Close closes the trace.
	Close(attrs ...Attr)

	// Error captures an event associated with this trace.
	Error(attrs ...Attr)

	// Warn captures an event associated with this trace.
	Warn(attrs ...Attr)

	// Info captures an event associated with this trace.
	Info(attrs ...Attr)

	// Debug captures an event associated with this trace.
	Debug(attrs ...Attr)

	// Log an event associated with this trace.
	Log(level Level, attrs ...Attr)

	// Assert that the conditions represented by the Attrs hold. If not, log
	// an event associated with the trace at the AssertionViolated level;
	// otherwise, log an event at the specified level.
	Assert(level Level, attrs ...Attr)

	// Count updates a counter associated with this trace. Counters are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Count(key Metric, delta int64)

	// Gauge updates a gauge associated with this trace. Gauges are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Gauge(key Metric, value int64)

	// Duration updates a timespan associated with this trace. Counters are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Duration(key Metric, d time.Duration)

	// Histogram updates a log-linear histogram associated with this trace. Histograms are
	// flushed when the span is closed. If telemetry is not configured, this operation
	// is a nop.
	Histogram(key Metric, sample int64)
}

var _ = Trace(&noptraceimpl{})
var _ = Trace(&traceimpl{})

type closer interface {
	close(Handler, []Attr)
}

var _ = closer(&noptraceimpl{})
var _ = closer(&traceimpl{})

type noptraceimpl struct{}

// singleton
var noptrace = &noptraceimpl{}

func (*noptraceimpl) ID() ksuid.KSUID                      { return ksuid.Nil }
func (*noptraceimpl) Elapsed() time.Duration               { return time.Duration(0) }
func (*noptraceimpl) Close(attrs ...Attr)                  {}
func (*noptraceimpl) Error(attrs ...Attr)                  {}
func (*noptraceimpl) Warn(attrs ...Attr)                   {}
func (*noptraceimpl) Info(attrs ...Attr)                   {}
func (*noptraceimpl) Debug(attrs ...Attr)                  {}
func (*noptraceimpl) Log(level Level, attrs ...Attr)       {}
func (*noptraceimpl) Assert(level Level, attrs ...Attr)    {}
func (*noptraceimpl) Count(key Metric, delta int64)        {}
func (*noptraceimpl) Gauge(key Metric, value int64)        {}
func (*noptraceimpl) Duration(key Metric, d time.Duration) {}
func (*noptraceimpl) Histogram(key Metric, sample int64)   {}
func (*noptraceimpl) close(Handler, []Attr)                {}

type traceimpl struct {
	pkg     *pkgimpl
	tp      *tracepoint
	id      ksuid.KSUID
	then    time.Time
	tpAttrs []Attr

	counts    map[Metric]int64
	gauges    map[Metric]int64
	durations map[Metric]time.Duration
}

func (tr *traceimpl) ID() ksuid.KSUID {
	return tr.id
}

func (tr *traceimpl) Elapsed() time.Duration {
	return time.Since(tr.then)
}

func (tr *traceimpl) Close(attrs ...Attr) {
	tr.pkg.close(tr.tp, tr, attrs)
}

func (tr *traceimpl) Error(attrs ...Attr) {
	tr.pkg.log(tr, 2, ErrorLevel, attrs)
}

func (tr *traceimpl) Warn(attrs ...Attr) {
	tr.pkg.log(tr, 2, WarnLevel, attrs)
}

func (tr *traceimpl) Info(attrs ...Attr) {
	tr.pkg.log(tr, 2, InfoLevel, attrs)
}

func (tr *traceimpl) Debug(attrs ...Attr) {
	tr.pkg.log(tr, 2, DebugLevel, attrs)
}

func (tr *traceimpl) Log(level Level, attrs ...Attr) {
	tr.pkg.log(tr, 2, level, attrs)
}

func (tr *traceimpl) Assert(level Level, attrs ...Attr) {
	for _, a := range attrs {
		if !a.Condition() {
			level = AssertionViolatedLevel
			break
		}
	}
	tr.pkg.log(tr, 2, level, attrs)
}

func (tr *traceimpl) Count(key Metric, delta int64) {
	if tr.counts == nil {
		tr.counts = make(map[Metric]int64)
	}
	val := delta
	if v, ok := tr.counts[key]; ok {
		val = val + v
	}
	tr.counts[key] = val
}

func (tr *traceimpl) Gauge(key Metric, value int64) {
	if tr.gauges == nil {
		tr.gauges = make(map[Metric]int64)
	}
	tr.gauges[key] = value
}

func (tr *traceimpl) Duration(key Metric, d time.Duration) {
	if tr.durations == nil {
		tr.durations = make(map[Metric]time.Duration)
	}
	tr.durations[key] = d
}

func (tr *traceimpl) Histogram(key Metric, sample int64) {
	// TODO: https://github.com/openhistogram/libcircllhist
}

func (tr *traceimpl) close(h Handler, attrs []Attr) {
	if src := tr.counts; src != nil {
		for m, v := range src {
			h.Count(m, v)
		}
	}
	if src := tr.gauges; src != nil {
		for m, v := range src {
			h.Gauge(m, v)
		}
	}
	if src := tr.durations; src != nil {
		for m, v := range src {
			h.Duration(m, v)
		}
	}
	h.Log(DebugLevel, [][]Attr{{
		Event("close trace"),
		String("id", tr.ID().String()),
		Duration("elapsed", tr.Elapsed()),
	}, attrs})
}
