package trace

import (
	"runtime"
	"time"

	"github.com/segmentio/ksuid"
)

type Trace interface {
	ID() ksuid.KSUID // ID returns the unique identifier for the trace
	Finish()         // Finish closes the trace

	// Title sets the title of the trace
	Title(title string) Trace

	// Error writes an event log associated with a span.
	Error(probe Probe, attrs ...Attr)

	// Warn writes an event log associated with a span.
	Warn(probe Probe, attrs ...Attr)

	// Info writes an event log associated with a span.
	Info(probe Probe, attrs ...Attr)

	// Debug writes an event log associated with a span.
	Debug(probe Probe, attrs ...Attr)

	// Log an event associated with a span. Log events are immediately flushed
	// to the Logger associated with the package. If telemetry is not configured,
	// this operation is a nop.
	Log(probe Probe, level Level, attrs ...Attr)

	// Assert that the conditions represented by the Attrs hold. If not, log
	// an event associated with the trace at the AssertionViolated level;
	// otherwise, log an event at the specified level.
	Assert(probe Probe, level Level, attrs ...Attr)

	// Count updates a counter associated with a span. Counters are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Count(probe Probe, delta int64) int64

	// Gauge updates a gauge associated with a span. Gauges are flushed when
	// the span is closed. If telemetry is not configured, this operation is a nop.
	Gauge(probe Probe, value int64)

	// Histogram updates a log-linear histogram associated with a span. Histograms are
	// flushed when the span is closed. If telemetry is not configured, this operation
	// is a nop.
	Histogram(probe Probe, sample int64)
}

// New creates and returns a Span.
func New(pkg Package, attrs ...Attr) Trace {
	tr := &traceimpl{
		pkg:    pkg,
		id:     ksuid.New(),
		title:  "anonymous",
		then:   time.Now(),
		attrs:  attrs,
		counts: make(map[Probe]int64),
		gauges: make(map[Probe]int64),
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

type traceimpl struct {
	pkg      Package
	id       ksuid.KSUID
	title    string
	then     time.Time
	duration time.Duration
	attrs    []Attr

	gid      uint64 // goroutine identifier
	file     string
	line     int
	funcName string

	counts map[Probe]int64
	gauges map[Probe]int64
}

func (tr *traceimpl) ID() ksuid.KSUID {
	return tr.id
}

func (tr *traceimpl) Title(title string) Trace {
	tr.title = title
	return tr
}

func (tr *traceimpl) Finish() {
	tr.duration = time.Since(tr.then)
	if h := tr.pkg.Handler(); h != nil {
		for p, val := range tr.counts {
			h.Count(tr, p, val)
		}
		for p, val := range tr.gauges {
			h.Gauge(tr, p, val)
		}
	}
}

func (tr *traceimpl) Error(probe Probe, attrs ...Attr) {
	tr.log(2, probe, ErrorLevel, attrs...)
}

func (tr *traceimpl) Warn(probe Probe, attrs ...Attr) {
	tr.log(2, probe, WarnLevel, attrs...)
}

func (tr *traceimpl) Info(probe Probe, attrs ...Attr) {
	tr.log(2, probe, InfoLevel, attrs...)
}

func (tr *traceimpl) Debug(probe Probe, attrs ...Attr) {
	tr.log(2, probe, DebugLevel, attrs...)
}

func (tr *traceimpl) Log(probe Probe, level Level, attrs ...Attr) {
	tr.log(2, probe, level, attrs...)
}

func (tr *traceimpl) log(skip int, probe Probe, level Level, attrs ...Attr) {
	if probe.Enabled(level) {
		if h := tr.pkg.Handler(); h != nil && h.Enabled(level) {
			var file string
			var line int
			var gid uint64
			if tr.pkg.CaptureSourceInfo() {
				_, file, line, _ = runtime.Caller(skip)
				gid = __caution__GetGoroutineID()
			}

			evt := NewEventLog(time.Now(), level, probe.String(), file, line, gid)
			for _, a := range attrs {
				evt.AddAttr(a)
			}
			for _, b := range tr.attrs {
				evt.AddAttr(b)
			}
			h.Log(tr, evt)
		}
	}
}

func (tr *traceimpl) Assert(probe Probe, level Level, attrs ...Attr) {
	for _, a := range attrs {
		if !a.Condition() {
			level = AssertionViolatedLevel
			break
		}
	}
	tr.log(2, probe, level, attrs...)
}

func (tr *traceimpl) Count(probe Probe, delta int64) (val int64) {
	val = delta
	if v, ok := tr.counts[probe]; ok {
		val = val + v
	}
	tr.counts[probe] = val
	return
}

func (tr *traceimpl) Gauge(probe Probe, value int64) {
	tr.gauges[probe] = value
}

func (tr *traceimpl) Histogram(probe Probe, sample int64) {
	// TODO: https://github.com/openhistogram/libcircllhist
}
