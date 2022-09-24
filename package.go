package trace

import (
	"runtime"
	"sync"
)

type PackageFlags int16

const (
	CaptureGoroutineID PackageFlags = 1 << iota
	CaptureSourceInfo  PackageFlags = 1 << iota
)

// Package represents the integration point with third-party packages that
// implement telemetry.  Libraries and other packages that support telemetry
// should export a top-level function that returns a package-local instance
// that implements this interface.
type Package interface {
	Flags() PackageFlags

	// Connect adds a Handler to a tracepoint.
	Connect(Handler, Point)

	// Disconnect removes a Handler from a tracepoint.
	Disconnect(Handler, Point)

	// Enabled returns whether any Handlers are connected to a tracepoint.
	Enabled(Point) bool
}

// NewPackage returns a new Package.
func NewPackage(flags PackageFlags) Package {
	return &pkgimpl{
		flags:    flags,
		mu:       sync.RWMutex{},
		handlers: make(map[Point]map[Handler]struct{}),
	}
}

type pkgimpl struct {
	flags    PackageFlags
	mu       sync.RWMutex
	handlers map[Point]map[Handler]struct{}
}

func (pkg *pkgimpl) Flags() PackageFlags {
	return pkg.flags
}

func (pkg *pkgimpl) Connect(h Handler, tracepoint Point) {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()
	m, ok := pkg.handlers[tracepoint]
	if !ok {
		m = make(map[Handler]struct{})
		pkg.handlers[tracepoint] = m
	}
	m[h] = struct{}{}
}

func (pkg *pkgimpl) Disconnect(h Handler, tracepoint Point) {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()
	if m, ok := pkg.handlers[tracepoint]; ok {
		delete(m, h)
	}
}

func (pkg *pkgimpl) Enabled(tracepoint Point) bool {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()
	if m, ok := pkg.handlers[tracepoint]; ok {
		return len(m) > 0
	}
	return false
}

func (pkg *pkgimpl) log(tr *traceimpl, skip int, level Level, attrs []Attr) {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()
	if m, ok := pkg.handlers[tr.tp]; ok {
		for h := range m {
			if h.Enabled(level) {
				if (pkg.flags & CaptureGoroutineID) == CaptureGoroutineID {
					if attrs == nil {
						attrs = make([]Attr, 0)
					}
					gid := __caution__GetGoroutineID()
					attrs = append(attrs, Uint64("gid", gid))
				}
				if (pkg.flags & CaptureSourceInfo) == CaptureSourceInfo {
					if attrs == nil {
						attrs = make([]Attr, 0)
					}
					_, file, line, _ := runtime.Caller(skip)
					attrs = append(attrs, String("file", file), Int("line", line))
				}

				arr := [][]Attr{tr.tpAttrs}
				if attrs != nil {
					arr = append(arr, attrs)
				}

				h.Log(level, arr)
			}
		}
	}
}

func (pkg *pkgimpl) begin(tr *traceimpl) {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()
	if m, ok := pkg.handlers[tr.tp]; ok {
		for h := range m {
			h.Log(DebugLevel, [][]Attr{{
				Event("begin trace"),
				String("id", tr.ID().String()),
			}})
		}
	}
}

func (pkg *pkgimpl) close(tp Point, tr Trace, attrs []Attr) {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()
	if m, ok := pkg.handlers[tp]; ok {
		if f, ok := tr.(closer); ok {
			for h := range m {
				f.close(h, attrs)
			}
		}
	}
}
