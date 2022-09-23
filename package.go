package trace

import (
	"errors"
)

// Package represents the integration point with third-party packages that
// implement telemetry.  Libraries and other packages that support telemetry
// should export a top-level function that returns a package-local instance
// that implements this interface.
type Package interface {
	Bind(*Router)
	Handler() Handler
	CaptureSourceInfo() bool

	// Skip configures the probe to be disabled for traces associated with this package.
	Skip(Probe)

	// Enabled returns whether or not the probe is enabled.
	Enabled(Probe) bool

	// Registry returns an application-specific mapping of probes to identifiers.
	Registry(app string) Registry
}

type pkgimpl struct {
	r          *Router
	h          Handler
	skipset    map[Probe]struct{}
	registries map[string]Registry
}

// NewPackage returns a new Package that delegates to a Router.
func NewPackage() Package {
	return &pkgimpl{
		skipset:    make(map[Probe]struct{}),
		registries: make(map[string]Registry),
	}
}

// NewPackageWithHandler returns a new Package that uses its own Handler.
func NewPackageWithHandler(h Handler) Package {
	return &pkgimpl{
		h:          h,
		skipset:    make(map[Probe]struct{}),
		registries: make(map[string]Registry),
	}
}

func (p *pkgimpl) Bind(r *Router) {
	p.r = r
}

var ErrMustUsePackage = errors.New("package not configured")

func (p *pkgimpl) Handler() Handler {
	switch {
	case p.h != nil:
		return p.h
	case p.r != nil:
		return p.r.Handler()
	default:
		panic(ErrMustUsePackage)
	}
}

func (p *pkgimpl) CaptureSourceInfo() bool {
	if p.r != nil {
		return p.r.SourceInfo()
	}
	return false
}

// Skip configures the probe to be disabled for traces associated with this package.
func (pkg *pkgimpl) Skip(p Probe) {
	pkg.skipset[p] = struct{}{}
}

func (pkg *pkgimpl) Enabled(p Probe) bool {
	_, ok := pkg.skipset[p]
	return !ok
}

func (pkg *pkgimpl) Registry(app string) Registry {
	a, ok := pkg.registries[app]
	if !ok {
		a = NewRegistry()
		pkg.registries[app] = a
	}
	return a
}
