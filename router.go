package trace

import (
	"sync"
)

// Router steers telemetry from Packages into backends.
type Router struct {
	sourceInfo bool
	handler    Handler
	pkgmu      sync.Mutex
	set        map[Package]struct{}
}

func NewRouter(h Handler, sourceInfo bool) *Router {
	return &Router{
		sourceInfo: sourceInfo,
		handler:    h,
		pkgmu:      sync.Mutex{},
		set:        make(map[Package]struct{}),
	}
}

// Use registers a package.
func (r *Router) Use(pkg Package) {
	r.pkgmu.Lock()
	defer r.pkgmu.Unlock()
	if _, ok := r.set[pkg]; !ok {
		r.set[pkg] = struct{}{}
	}
	pkg.Bind(r)
}

func (r *Router) Handler() Handler {
	return r.handler
}

func (r *Router) SourceInfo() bool {
	return r.sourceInfo
}
