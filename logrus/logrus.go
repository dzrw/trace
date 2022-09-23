package logrus

import (
	"fmt"
	"time"

	"github.com/dzrw/trace"
	log "github.com/sirupsen/logrus"
)

type LogrusHandler struct {
	minLevel trace.Level
	l        *log.Logger
}

var _ = trace.Handler(&LogrusHandler{})

func NewLogrusHandler(logger *log.Logger) trace.Handler {
	minLevel := ToTraceLevel(logger.GetLevel())
	return &LogrusHandler{minLevel, logger}
}

func (h *LogrusHandler) Enabled(l trace.Level) bool {
	return l <= h.minLevel
}

func (h *LogrusHandler) Log(tr trace.Trace, evt *trace.EventLog) error {
	if evt.Time().IsZero() || evt.Level() == 0 {
		return nil
	}

	f := log.NewEntry(h.l)
	f.WithTime(evt.Time())
	f.WithField("trace", tr.ID().String())
	if gid := evt.Goroutine(); gid > 0 {
		f.WithField("gid", gid)
	}
	if file, line := evt.SourceLine(); file != "" {
		f.WithField("file", fmt.Sprint(file, ":", line))
	}
	for _, attr := range evt.Attrs() {
		f.WithField(attr.Format())
	}

	f.Log(ToLogrusLevel(evt.Level()), evt.Message())
	return nil
}

func (h *LogrusHandler) Count(tr trace.Trace, p trace.Probe, delta int64, attrs ...trace.Attr) error {
	evt := trace.NewEventLog(time.Now(), trace.InfoLevel, p.String(), "", 0, 0)
	evt.AddAttr(trace.String("count", fmt.Sprint(delta)))
	for _, a := range attrs {
		evt.AddAttr(a)
	}
	return h.Log(tr, evt)
}

func (h *LogrusHandler) Gauge(tr trace.Trace, p trace.Probe, value int64, attrs ...trace.Attr) error {
	evt := trace.NewEventLog(time.Now(), trace.InfoLevel, p.String(), "", 0, 0)
	evt.AddAttr(trace.String("gauge", fmt.Sprint(value)))
	for _, a := range attrs {
		evt.AddAttr(a)
	}
	return h.Log(tr, evt)
}

func ToTraceLevel(l log.Level) trace.Level {
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

func ToLogrusLevel(l trace.Level) log.Level {
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
