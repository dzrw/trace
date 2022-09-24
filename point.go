package trace

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/segmentio/ksuid"
)

type Point interface {
	fmt.Stringer
	Capture(...Attr) Trace // New returns a new Trace spawned from this tracepoint.
	Enabled() bool         // Enabled returns whether or not this tracepoint is enabled.
}

func NewPoint(pkg Package, name string, attrs ...Attr) Point {
	tp := &tracepoint{
		pkg:   pkg,
		name:  name,
		attrs: attrs,
	}
	return tp
}

type tracepoint struct {
	pkg   Package
	name  string
	attrs []Attr
}

func (tp *tracepoint) Capture(attrs ...Attr) Trace {
	if !tp.pkg.Enabled(tp) {
		return noptrace
	}

	flags := tp.pkg.Flags()
	if (flags & CaptureGoroutineID) == CaptureGoroutineID {
		if gid := __caution__GetGoroutineID(); gid > 0 {
			attrs = append(attrs, Uint64("gid", gid))
		}
	}

	if (flags & CaptureSourceInfo) == CaptureSourceInfo {
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
		pkg:     tp.pkg.(*pkgimpl),
		tp:      tp,
		id:      ksuid.New(),
		then:    time.Now(),
		tpAttrs: attrs,
	}

	tr.pkg.begin(tr)
	return tr
}

func (tp *tracepoint) Enabled() bool {
	return tp.pkg.Enabled(tp)
}

func (tp *tracepoint) String() string {
	return tp.name
}
