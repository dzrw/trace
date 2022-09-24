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
}

func NewHandler(logger *log.Logger) *LogrusHandler {
	minLevel := makeTraceLevel(logger.GetLevel())
	return &LogrusHandler{minLevel, logger}
}

func (h *LogrusHandler) Capabilities() trace.HandlerCaps {
	return trace.SupportsLogs | trace.SupportsMetrics | trace.SupportsTraces
}

func (h *LogrusHandler) Enabled(l trace.Level) bool {
	return l <= h.minLevel
}

func (h *LogrusHandler) Log(l trace.Level, attrs [][]trace.Attr) error {
	if l == 0 {
		return nil
	}

	m := make(log.Fields)
	format3(m, attrs)
	f := log.NewEntry(h.logger).WithFields(m)
	f.Log(makeLogrusLevel(l))
	return nil
}

func (h *LogrusHandler) Count(m trace.Metric, delta int64) error {
	f := make(log.Fields)
	format1(f, trace.Int64(string(m), delta))
	e := log.NewEntry(h.logger).WithFields(f)
	e.Log(log.InfoLevel)
	return nil
}

func (h *LogrusHandler) Gauge(m trace.Metric, value int64) error {
	f := make(log.Fields)
	format1(f, trace.Int64(string(m), value))
	e := log.NewEntry(h.logger).WithFields(f)
	e.Log(log.InfoLevel)
	return nil
}

func (h *LogrusHandler) Duration(m trace.Metric, d time.Duration) error {
	f := make(log.Fields)
	format1(f, trace.Duration(string(m), d))
	e := log.NewEntry(h.logger).WithFields(f)
	e.Log(log.InfoLevel)
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
		format2(m, arr)
	}
}

func format2(m log.Fields, attrs []trace.Attr) {
	for _, a := range attrs {
		format1(m, a)
	}
}

func format1(m log.Fields, a trace.Attr) {
	k, v := a.Format()
	m[k] = trace.Quote(v)
}
