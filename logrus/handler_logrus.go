package logrus

import (
	"time"

	"github.com/dzrw/trace"
	log "github.com/sirupsen/logrus"
)

var _ = trace.Handler(&LogrusHandler{})

type LogrusHandler struct {
	minLevel trace.Level
	logger   *log.Logger
	flags    trace.HandlerFlags
	reg      trace.Registry
}

func New(logger *log.Logger, includeSourceInfo, includeGoroutineID bool, reg trace.Registry) *LogrusHandler {
	var flags trace.HandlerFlags
	if includeSourceInfo {
		flags |= trace.FlagSourceInfo
	}
	if includeGoroutineID {
		flags |= trace.FlagGoroutineID
	}
	if reg == nil {
		reg = trace.NewRegistry()
	}
	minLevel := makeTraceLevel(logger.GetLevel())
	return &LogrusHandler{minLevel, logger, flags, reg}
}

func (h *LogrusHandler) Flags() trace.HandlerFlags {
	return h.flags
}

func (h *LogrusHandler) Enabled(l trace.Level) bool {
	return l <= h.minLevel
}

func (h *LogrusHandler) TraceCreated(tr trace.Trace, attrs []trace.Attr) {
	h.Log(tr, trace.DebugLevel, []trace.Attr{
		trace.Event("trace created"),
	}, attrs)
}

func (h *LogrusHandler) TraceFinished(tr trace.Trace, attrs []trace.Attr) {
	h.Log(tr, trace.DebugLevel, []trace.Attr{
		trace.Event("trace finished"),
		trace.Duration("elapsed", tr.Elapsed()),
	}, attrs)
}

func (h *LogrusHandler) Count(tp trace.Tracepoint, delta int64) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		f := make(log.Fields)
		format2(f,
			trace.String("site", str),
			trace.Int64("count", delta))
		e := log.NewEntry(h.logger).WithFields(f)
		e.Log(log.InfoLevel)
	}
	return nil
}

func (h *LogrusHandler) Gauge(tp trace.Tracepoint, value int64) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		f := make(log.Fields)
		format2(f,
			trace.String("site", str),
			trace.Int64("gauge", value))
		e := log.NewEntry(h.logger).WithFields(f)
		e.Log(log.InfoLevel)
	}
	return nil
}

func (h *LogrusHandler) Duration(tp trace.Tracepoint, d time.Duration) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		f := make(log.Fields)
		format2(f,
			trace.String("site", str),
			trace.Duration("duration", d))
		e := log.NewEntry(h.logger).WithFields(f)
		e.Log(log.InfoLevel)
	}
	return nil
}

func (h *LogrusHandler) Histogram(tp trace.Tracepoint, sample int64) error {
	return nil // TODO
}

func (h *LogrusHandler) Log(tr trace.Trace, l trace.Level, attrs ...[]trace.Attr) error {
	if l == 0 {
		return nil
	}

	if site, ok := h.reg.IdentifierFor(tr.Site()); ok {
		f := make(log.Fields)
		format2(f,
			trace.String("site", site),
			trace.Uint64("trace", tr.ID()),
		)
		format3(f, attrs)
		e := log.NewEntry(h.logger).WithFields(f)
		e.Log(makeLogrusLevel(l))
	}

	return nil
}

func makeTraceLevel(l log.Level) trace.Level {
	switch l {
	case log.ErrorLevel:
		return trace.ErrorLevel
	case log.WarnLevel:
		return trace.WarnLevel
	case log.InfoLevel:
		return trace.InfoLevel
	default:
		return trace.DebugLevel
	}
}

func makeLogrusLevel(l trace.Level) log.Level {
	switch {
	case l <= trace.ErrorLevel:
		return log.ErrorLevel
	case l <= trace.WarnLevel:
		return log.WarnLevel
	case l <= trace.InfoLevel:
		return log.InfoLevel
	default:
		return log.DebugLevel
	}
}

func format3(m log.Fields, attrs [][]trace.Attr) {
	for _, arr := range attrs {
		format2(m, arr...)
	}
}

func format2(m log.Fields, attrs ...trace.Attr) {
	for _, a := range attrs {
		format1(m, a)
	}
}

func format1(m log.Fields, a trace.Attr) {
	k, v := a.Format()
	m[k] = v
}
