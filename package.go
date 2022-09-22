package trace

// Package represents the integration point with third-party packages that
// implement telemetry.  Libraries and other packages that support telemetry
// should export a top-level function that returns a package-local instance
// that implements this interface.
type Package interface {
	Bind(*Router)
	Handler() Handler
	CaptureSourceInfo() bool
}

type pkgimpl struct {
	r *Router
	h Handler
}

// NewPackage returns a new Package that delegates to a Router.
func NewPackage() Package {
	return &pkgimpl{}
}

// NewPackageWithHandler returns a new Package that uses its own Handler.
func NewPackageWithHandler(h Handler) Package {
	return &pkgimpl{h: h}
}

func (p *pkgimpl) Bind(r *Router) {
	p.r = r
}

func (p *pkgimpl) Handler() Handler {
	switch {
	case p.h != nil:
		return p.h
	case p.r != nil:
		return p.r.Handler()
	default:
		panic("package not configured")
	}
}

func (p *pkgimpl) CaptureSourceInfo() bool {
	if p.r != nil {
		return p.r.SourceInfo()
	}
	return false
}
