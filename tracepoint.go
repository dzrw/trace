package trace

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/segmentio/ksuid"
)

/*
Usage:
	var SiteFoo = trace.Site() // exports a tracepoint

	func foo() {
		SiteFoo.Count(1) // counter

		tr := SiteFoo.Trace()
		defer tr.Close()

		tr.Debug("foo") // log
	}
*/

type Tracepoint interface {
	ID() ksuid.KSUID // ID returns the unique identifier of this tracepoint.

	Install(Handler)          // Install a Handler into this tracepoint.
	Uninstall()               // Uninstall the Handler from this tracepoint.
	Handler() (Handler, bool) // Handler returns the current Handler, if any.

	Trace(...Attr) Trace // Trace originates a new Trace from this tracepoint.

	Count(delta int64)        // Count captures a delta from this tracepoint.
	Gauge(value int64)        // Gauge captures a value from this tracepoint.
	Duration(d time.Duration) // Duration captures a duration from this tracepoint.
	Histogram(sample int64)   // Histogram captures a sample from this tracepoint.
}

func Site() Tracepoint {
	tp := &tracepoint{
		id: ksuid.New(),
	}
	return tp
}

type tracepoint struct {
	id   ksuid.KSUID
	next uint64
	h    Handler
}

func (tp *tracepoint) ID() ksuid.KSUID {
	return tp.id
}

func (tp *tracepoint) Install(h Handler) {
	tp.h = h
}

func (tp *tracepoint) Uninstall() {
	tp.h = nil
}

func (tp *tracepoint) Handler() (h Handler, ok bool) {
	return tp.h, tp.h != nil
}

func (tp *tracepoint) Trace(attrs ...Attr) Trace {
	if h, ok := tp.Handler(); ok {
		flags := h.Flags()
		if (flags & FlagGoroutineID) == FlagGoroutineID {
			if gid := __caution__GetGoroutineID(); gid > 0 {
				attrs = append(attrs, Uint64("gid", gid))
			}
		}

		if (flags & FlagSourceInfo) == FlagSourceInfo {
			if pc, file, line, ok := runtime.Caller(1); ok {
				fileAttr := String("file", fmt.Sprint(file, ":", strconv.Itoa(line)))
				attrs = append(attrs, fileAttr)
				if f := runtime.FuncForPC(pc); f != nil {
					funcAttr := String("func", f.Name())
					attrs = append(attrs, funcAttr)
				}
			}
		}

		tr := &traceimpl{
			tp:    tp,
			id:    tp.next,
			then:  time.Now(),
			attrs: attrs,
		}

		tp.next++
		h.TraceCreated(tr, tr.attrs)
		return tr
	}

	return noptrace

}

func (tp *tracepoint) Count(delta int64) {
	if h, ok := tp.Handler(); ok {
		h.Count(tp, delta)
	}
}

func (tp *tracepoint) Gauge(value int64) {
	if h, ok := tp.Handler(); ok {
		h.Gauge(tp, value)
	}
}

func (tp *tracepoint) Duration(d time.Duration) {
	if h, ok := tp.Handler(); ok {
		h.Duration(tp, d)
	}
}

func (tp *tracepoint) Histogram(sample int64) {
	if h, ok := tp.Handler(); ok {
		_ = h
		// TODO: https://github.com/openhistogram/libcircllhist
	}
}

func (tp *tracepoint) log(tr Trace, skip int, level Level, attrs []Attr) {
	if h, ok := tp.Handler(); ok {
		if h.Enabled(level) {
			flags := h.Flags()

			if (flags & FlagGoroutineID) == FlagGoroutineID {
				if attrs == nil {
					attrs = make([]Attr, 0)
				}
				gid := __caution__GetGoroutineID()
				attrs = append(attrs, Uint64("gid", gid))
			}

			if (flags & FlagSourceInfo) == FlagSourceInfo {
				if attrs == nil {
					attrs = make([]Attr, 0)
				}
				_, file, line, _ := runtime.Caller(skip)
				attrs = append(attrs, String("file", file), Int("line", line))
			}

			h.Log(tr, level, attrs)
		}
	}
}

func (tp *tracepoint) finishTrace(tr Trace, attrs []Attr) {
	if h, ok := tp.Handler(); ok {
		h.TraceFinished(tr, attrs)
	}
}
