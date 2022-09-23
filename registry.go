package trace

type Registry interface {
	Count() int
	Define(Probe, string)
	IsDefined(Probe) bool
	IdentifierFor(Probe) (string, bool)
	ProbeFor(string) (Probe, bool)
}

type registry struct {
	u map[Probe]string
	v map[string]Probe
}

func NewRegistry() Registry {
	return &registry{
		u: make(map[Probe]string),
		v: make(map[string]Probe),
	}
}

func (m *registry) Count() int {
	return len(m.u)
}

func (m *registry) Define(p Probe, id string) {
	m.u[p] = id
	m.v[id] = p
}

func (m *registry) IsDefined(p Probe) bool {
	_, ok := m.u[p]
	return ok
}

func (m *registry) IdentifierFor(p Probe) (id string, ok bool) {
	id, ok = m.u[p]
	return
}
func (m *registry) ProbeFor(id string) (p Probe, ok bool) {
	p, ok = m.v[id]
	return
}
