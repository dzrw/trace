package trace

// Registry provides a one-to-one mapping from Tracepoints to
// application-specific identifiers. Applications should use this type to
// help do stuff.
type Registry interface {
	Count() int
	Define(Tracepoint, string)
	Undefine(Tracepoint)
	IsDefined(Tracepoint) bool
	IdentifierFor(Tracepoint) (string, bool)
	TracepointFor(string) (Tracepoint, bool)
}

type registry struct {
	u map[Tracepoint]string
	v map[string]Tracepoint
}

func NewRegistry() Registry {
	return &registry{
		u: make(map[Tracepoint]string),
		v: make(map[string]Tracepoint),
	}
}

func (m *registry) Count() int {
	return len(m.u)
}

func (m *registry) Define(tp Tracepoint, id string) {
	m.u[tp] = id
	m.v[id] = tp
}

func (m *registry) Undefine(tp Tracepoint) {
	if id, ok := m.u[tp]; ok {
		delete(m.u, tp)
		delete(m.v, id)
	}
}

func (m *registry) IsDefined(tp Tracepoint) bool {
	_, ok := m.u[tp]
	return ok
}

func (m *registry) IdentifierFor(tp Tracepoint) (id string, ok bool) {
	id, ok = m.u[tp]
	return
}

func (m *registry) TracepointFor(id string) (tp Tracepoint, ok bool) {
	tp, ok = m.v[id]
	return
}
