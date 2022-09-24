package text

import (
	"io"
	"strings"
	"sync"
	"time"

	"github.com/dzrw/trace"
)

var _ = trace.Handler(&TextHandler{})

// TextHandler is a Handler that writes to an io.Writer.
type TextHandler struct {
	// Has unexported fields.
	w  io.Writer
	l  trace.Level
	mu sync.Mutex
}

// NewTextHandler creates a TextHandler that writes to w using the default
// options.
func NewHandler(w io.Writer, l trace.Level) *TextHandler {
	return &TextHandler{
		w:  w,
		l:  l,
		mu: sync.Mutex{},
	}
}

func (h *TextHandler) Capabilities() trace.HandlerCaps {
	return trace.SupportsLogs | trace.SupportsMetrics | trace.SupportsTraces
}

// Enabled reports whether this handler is accepting event logs at the given level.
func (h *TextHandler) Enabled(l trace.Level) bool {
	return l <= h.l
}

/*
Log formats its argument EventLog as a single line of space-separated
key=value items.

If the EventLog's time is zero, it is omitted. Otherwise, the key is "time"
and the value is output in RFC3339 format with millisecond precision.

If the EventLog's level is zero, it is omitted. Otherwise, the key is "level"
and the value of Level.String is output.

If the AddSource option is set and source information is available,
the key is "source" and the value is output as FILE:LINE.

The message's key "msg".

To modify these or other attributes, or remove them from the output,
use [LoggerOptions.ReplaceAttr].

Keys are written as unquoted strings. Values are written according to their
type:
  - Strings are quoted if they contain Unicode space characters or are over
    80 bytes long.
  - If a value implements [encoding.TextMarshaler], the result of
    MarshalText is used.
  - Otherwise, the result of fmt.Sprint is used.

Each call to Handle results in a single, mutex-protected call to
io.Writer.Write.
*/
func (h *TextHandler) Log(l trace.Level, attrs [][]trace.Attr) error {
	if l == 0 {
		return nil
	}

	sb := strings.Builder{}
	format3(&sb, attrs)
	return h.finish(&sb)
}

func (h *TextHandler) Count(m trace.Metric, delta int64) error {
	sb := strings.Builder{}
	format1(&sb, trace.Int64(string(m), delta))
	return h.finish(&sb)
}

func (h *TextHandler) Gauge(m trace.Metric, value int64) error {
	sb := strings.Builder{}
	format1(&sb, trace.Int64(string(m), value))
	return h.finish(&sb)
}

func (h *TextHandler) Duration(m trace.Metric, d time.Duration) error {
	sb := strings.Builder{}
	format1(&sb, trace.Duration(string(m), d))
	return h.finish(&sb)
}

func (h *TextHandler) finish(sb *strings.Builder) error {
	sb.WriteString("\n")

	buf := []byte(sb.String())
	if len(buf) > 0 {
		buf = buf[1:]
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf)
	return err
}

func format3(sb *strings.Builder, attrs [][]trace.Attr) {
	for _, arr := range attrs {
		format2(sb, arr)
	}
}

func format2(sb *strings.Builder, attrs []trace.Attr) {
	for _, a := range attrs {
		format1(sb, a)
	}
}

func format1(sb *strings.Builder, a trace.Attr) {
	k, v := a.Format()
	sb.WriteRune(' ')
	sb.WriteString(k)
	sb.WriteRune('=')
	sb.WriteString(trace.Quote(v))
}
