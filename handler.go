package trace

import (
	"encoding"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A Handler processes traces and associated event logs.
type Handler interface {
	// Enabled reports whether this handler is accepting event logs
	// at the given level.
	Enabled(Level) bool

	// Log handles an event associated with a trace.
	Log(Tracer, *EventLog) error

	// Count handles a counter associated with a trace.
	Count(Tracer, Probe, int64, ...Attr) error

	// Gauge handles a gauge associated with a trace.
	Gauge(Tracer, Probe, int64, ...Attr) error

	// With returns a new Handler whose attributes consist of
	// the receiver's attributes concatenated with the arguments.
	//With([]Attr) Handler
}

// TextHandler is a Handler that writes Traces and EventLogs to an io.Writer.
type TextHandler struct {
	// Has unexported fields.
	w  io.Writer
	l  Level
	mu sync.Mutex
}

// NewTextHandler creates a TextHandler that writes to w using the default
// options.
func NewTextHandler(w io.Writer, l Level) *TextHandler {
	return &TextHandler{
		w:  w,
		l:  l,
		mu: sync.Mutex{},
	}
}

// Enabled reports whether this handler is accepting event logs at the given level.
func (h *TextHandler) Enabled(l Level) bool {
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
func (h *TextHandler) Log(tr Tracer, evt *EventLog) error {
	const RFC3339Ms = "2006-01-02T15:04:05.999Z07:00"

	et := evt.Time()
	el := evt.Level()

	if et.IsZero() || el == 0 {
		return nil
	}

	sb := strings.Builder{}
	sb.WriteString(" time=")
	sb.WriteString(et.Format(RFC3339Ms))

	sb.WriteString(" trace=")
	sb.WriteString(tr.ID().String())

	if gid := evt.Goroutine(); gid > 0 {
		sb.WriteString(" [")
		sb.WriteString(fmt.Sprint(gid))
		sb.WriteString("]")
	}

	sb.WriteString(" level=")
	sb.WriteString(el.String())

	if file, line := evt.SourceLine(); file != "" {
		sb.WriteString(" file=")
		sb.WriteString(file)
		sb.WriteRune(':')
		sb.WriteString(strconv.Itoa(line))
	}

	sb.WriteString(" msg=")
	sb.WriteString(quote(evt.Message()))

	for _, attr := range evt.Attrs() {
		var str string
		switch attr.Kind() {
		case AnyKind:
			if v := attr.Value(); v == nil {
				continue
			} else if m, ok := attr.Value().(encoding.TextMarshaler); ok {
				text, err := m.MarshalText()
				if err != nil {
					str = err.Error()
				} else {
					str = string(text)
				}
			} else {
				str = fmt.Sprint(attr.Value())
			}
		case BoolKind:
			str = fmt.Sprint(attr.Bool())
		case DurationKind:
			str = fmt.Sprint(attr.Duration())
		case Float64Kind:
			str = fmt.Sprint(attr.Float64())
		case Int64Kind:
			str = fmt.Sprint(attr.Int64())
		case StringKind:
			str = quote(attr.String())
		case TimeKind:
			str = attr.Time().Format(RFC3339Ms)
		case Uint64Kind:
			str = fmt.Sprint(attr.Uint64())
		}

		sb.WriteRune(' ')
		sb.WriteString(attr.Key())
		sb.WriteRune('=')
		sb.WriteString(str)
	}

	sb.WriteString("\n")

	h.mu.Lock()
	defer h.mu.Unlock()
	p := []byte(sb.String())
	_, err := h.w.Write(p[1:])
	return err
}

func quote(s string) string {
	if len(s) > 80 || strings.ContainsAny(s, " \t") {
		return strings.Join([]string{"\"", s, "\""}, "")
	}
	return s
}

func (h *TextHandler) Count(tr Tracer, p Probe, v int64, attrs ...Attr) error {
	evt := NewEventLog(time.Now(), InfoLevel, p.String(), "", 0, 0)
	evt.AddAttr(String("count", fmt.Sprint(v)))
	for _, a := range attrs {
		evt.AddAttr(a)
	}
	return h.Log(tr, evt)
}

func (h *TextHandler) Gauge(tr Tracer, p Probe, v int64, attrs ...Attr) error {
	evt := NewEventLog(time.Now(), InfoLevel, p.String(), "", 0, 0)
	evt.AddAttr(String("gauge", fmt.Sprint(v)))
	for _, a := range attrs {
		evt.AddAttr(a)
	}
	return h.Log(tr, evt)
}

// With returns a new TextLogger whose attributes consists of h's attributes
// followed by attrs.
// func (h *TextLogger) With(attrs []Attr) Logger {
// 	return nil
// }
