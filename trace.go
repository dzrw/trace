package trace

import (
	"time"
)

type Trace interface {
	Site() Tracepoint       // Site returns the tracepoint that originated this trace.
	ID() uint64             // ID returns the unique identifier for the trace
	Elapsed() time.Duration // Elapsed returns the time elapsed since the trace started.

	// Close closes the trace.
	Close(attrs ...Attr)

	// Error captures an event associated with this trace.
	Error(event string, attrs ...Attr)

	// Warn captures an event associated with this trace.
	Warn(event string, attrs ...Attr)

	// Info captures an event associated with this trace.
	Info(event string, attrs ...Attr)

	// Debug captures an event associated with this trace.
	Debug(event string, attrs ...Attr)

	// Log an event associated with this trace.
	Log(level Level, attrs ...Attr)

	// Assert that the conditions represented by the Attrs hold. If not, log
	// an event associated with the trace at the AssertionViolated level;
	// otherwise, log an event at the specified level.
	Assert(level Level, attrs ...Attr)
}

var _ = Trace(&noptraceimpl{})
var _ = Trace(&traceimpl{})

type noptraceimpl struct{}

// singleton
var noptrace = &noptraceimpl{}

func (*noptraceimpl) Site() Tracepoint                  { return nil }
func (*noptraceimpl) ID() uint64                        { return 0 }
func (*noptraceimpl) Elapsed() time.Duration            { return time.Duration(0) }
func (*noptraceimpl) Close(attrs ...Attr)               {}
func (*noptraceimpl) Error(event string, attrs ...Attr) {}
func (*noptraceimpl) Warn(event string, attrs ...Attr)  {}
func (*noptraceimpl) Info(event string, attrs ...Attr)  {}
func (*noptraceimpl) Debug(event string, attrs ...Attr) {}
func (*noptraceimpl) Log(level Level, attrs ...Attr)    {}
func (*noptraceimpl) Assert(level Level, attrs ...Attr) {}

type traceimpl struct {
	tp    *tracepoint
	id    uint64
	then  time.Time
	attrs []Attr
}

func (tr *traceimpl) Site() Tracepoint {
	return tr.tp
}

func (tr *traceimpl) ID() uint64 {
	return tr.id
}

func (tr *traceimpl) Elapsed() time.Duration {
	return time.Since(tr.then)
}

func (tr *traceimpl) Close(attrs ...Attr) {
	tr.tp.finishTrace(tr, attrs)
}

func (tr *traceimpl) Error(event string, attrs ...Attr) {
	tr.tp.log(tr, 2, ErrorLevel, append(attrs, Event(event)))
}

func (tr *traceimpl) Warn(event string, attrs ...Attr) {
	tr.tp.log(tr, 2, WarnLevel, append(attrs, Event(event)))
}

func (tr *traceimpl) Info(event string, attrs ...Attr) {
	tr.tp.log(tr, 2, InfoLevel, append(attrs, Event(event)))
}

func (tr *traceimpl) Debug(event string, attrs ...Attr) {
	tr.tp.log(tr, 2, DebugLevel, append(attrs, Event(event)))
}

func (tr *traceimpl) Log(level Level, attrs ...Attr) {
	tr.tp.log(tr, 2, level, attrs)
}

func (tr *traceimpl) Assert(level Level, attrs ...Attr) {
	for _, a := range attrs {
		if !a.Condition() {
			level = AssertionViolatedLevel
			break
		}
	}
	tr.tp.log(tr, 2, level, attrs)
}
