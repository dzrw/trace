package trace

import (
	"sync"
)

type ConfigFunc func(Package, Handler, bool)

// Router configures Packages and tracepoints.
type Router interface {
	Use(pkg Package, cfg ConfigFunc)
}

type routerimpl struct {
	pkgmu  sync.RWMutex
	pkgset map[Package]ConfigFunc
}

func NewRouter() Router {
	return &routerimpl{
		pkgmu:  sync.RWMutex{},
		pkgset: make(map[Package]ConfigFunc),
	}
}

// Use registers a package.
func (r *routerimpl) Use(pkg Package, cfg ConfigFunc) {
	r.pkgmu.Lock()
	defer r.pkgmu.Unlock()
	if _, ok := r.pkgset[pkg]; !ok {
		r.pkgset[pkg] = cfg
	}
}

func (r *routerimpl) Connect(h Handler) {
	r.configure(h, true)
}

func (r *routerimpl) Disconnect(h Handler) {
	r.configure(h, false)
}

func (r *routerimpl) configure(h Handler, connect bool) {
	r.pkgmu.RLock()
	defer r.pkgmu.RUnlock()
	for pkg, cfg := range r.pkgset {
		cfg(pkg, h, connect)
	}
}
