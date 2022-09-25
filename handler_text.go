package trace

import (
	"io"
	"strings"
	"sync"
	"time"
)

var _ = Handler(&TextHandler{})

// TextHandler is a Handler that writes to an io.Writer.
type TextHandler struct {
	flags HandlerFlags
	l     Level
	reg   Registry
	mu    sync.Mutex
	w     io.Writer
}

// NewTextHandler creates a TextHandler that writes to w using the default
// options.
func NewTextHandler(w io.Writer, l Level, includeSourceInfo, includeGoroutineID bool, reg Registry) *TextHandler {
	var flags HandlerFlags
	if includeSourceInfo {
		flags |= FlagSourceInfo
	}
	if includeGoroutineID {
		flags |= FlagGoroutineID
	}
	if reg == nil {
		reg = NewRegistry()
	}
	return &TextHandler{
		flags: flags,
		l:     l,
		reg:   reg,
		mu:    sync.Mutex{},
		w:     w,
	}
}

func (h *TextHandler) Flags() HandlerFlags {
	return h.flags
}

func (h *TextHandler) Enabled(l Level) bool {
	return l <= h.l
}

func (h *TextHandler) TraceCreated(tr Trace, attrs []Attr) {
	h.Log(tr, DebugLevel, []Attr{
		Event("trace created"),
	}, attrs)
}

func (h *TextHandler) TraceFinished(tr Trace, attrs []Attr) {
	h.Log(tr, DebugLevel, []Attr{
		Event("trace finished"),
		Duration("elapsed", tr.Elapsed()),
	}, attrs)
}

func (h *TextHandler) Count(tp Tracepoint, delta int64) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		sb := strings.Builder{}
		format2(&sb,
			String("site", str),
			Int64("count", delta))
		return h.finish(&sb)
	}
	return nil
}

func (h *TextHandler) Gauge(tp Tracepoint, value int64) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		sb := strings.Builder{}
		format2(&sb,
			String("site", str),
			Int64("gauge", value))
		return h.finish(&sb)
	}
	return nil
}

func (h *TextHandler) Duration(tp Tracepoint, d time.Duration) error {
	if str, ok := h.reg.IdentifierFor(tp); ok {
		sb := strings.Builder{}
		format2(&sb,
			String("site", str),
			Duration("duration", d))
		return h.finish(&sb)
	}
	return nil
}

func (h *TextHandler) Histogram(tp Tracepoint, sample int64) error {
	return nil // TODO
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
func (h *TextHandler) Log(tr Trace, l Level, attrs ...[]Attr) error {
	if l == 0 {
		return nil
	}

	if site, ok := h.reg.IdentifierFor(tr.Site()); ok {
		sb := strings.Builder{}
		format2(&sb,
			String("site", site),
			Uint64("trace", tr.ID()),
		)
		format3(&sb, attrs)
		return h.finish(&sb)
	}

	return nil
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

func format3(sb *strings.Builder, attrs [][]Attr) {
	for _, arr := range attrs {
		format2(sb, arr...)
	}
}

func format2(sb *strings.Builder, attrs ...Attr) {
	for _, a := range attrs {
		format1(sb, a)
	}
}

func format1(sb *strings.Builder, a Attr) {
	k, v := a.Format()
	sb.WriteRune(' ')
	sb.WriteString(k)
	sb.WriteRune('=')
	sb.WriteString(quote(v))
}

func quote(s string) string {
	if len(s) > 80 || strings.ContainsAny(s, " \t") {
		return strings.Join([]string{"\"", s, "\""}, "")
	}
	return s
}
