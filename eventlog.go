package trace

import "time"

// An EventLog holds information about a log event.
type EventLog struct {
	// Has unexported fields.
	t     time.Time
	level Level
	msg   string
	gid   uint64
	file  string
	line  int
	attrs []Attr
}

// NewEventLog creates a new EventLog from the given arguments. Use EventLog.AddAttr
// to add attributes to the EventLog.

// NewEventLog is intended for logging APIs that want to support a Handler as a
// backend. Most users won't need it.
func NewEventLog(t time.Time, level Level, msg string, file string, line int) *EventLog {
	u := &EventLog{
		t:     t,
		level: level,
		msg:   msg,
		file:  file,
		line:  line,
	}

	if file != "" {
		u.gid = __caution__GetGoroutineID()
	}

	return u
}

// AddAttr appends a to the list of r's attributes. It does not check for
// duplicate keys.
func (r *EventLog) AddAttr(a Attr) {
	r.attrs = append(r.attrs, a)
}

// Attr returns the i'th Attr in r.
func (r *EventLog) Attr(i int) Attr {
	return r.attrs[i]
}

// Attrs returns a copy of the sequence of Attrs in r.
func (r *EventLog) Attrs() []Attr {
	dst := make([]Attr, len(r.attrs))
	copy(dst, r.attrs)
	return dst
}

// Level returns the level of the log event.
func (r *EventLog) Level() Level {
	return r.level
}

// Message returns the log message.
func (r *EventLog) Message() string {
	return r.msg
}

// NumAttrs returns the number of Attrs in r.
func (r *EventLog) NumAttrs() int {
	return len(r.attrs)
}

// SourceLine returns the file and line of the log event. If the EventLog
// was created without the necessary information, or if the location is
// unavailable, it returns ("", 0).
func (r *EventLog) SourceLine() (file string, line int) {
	return r.file, r.line
}

// Goroutine returns the identifier of the goroutine that created the event.
func (r *EventLog) Goroutine() uint64 {
	return r.gid
}

// Time returns the time of the log event.
func (r *EventLog) Time() time.Time {
	return r.t
}
