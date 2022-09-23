package trace

import (
	"errors"
	"fmt"
)

var ErrUnresolved = errUnresolved
var errUnresolved = errors.New("missing or unresolvable probe")

type Probe interface {
	fmt.Stringer
	Enabled(l Level) bool
}

func NewProbe(text string, minLevel Level) Probe {
	return &probe{text, minLevel}
}

type probe struct {
	text     string
	minLevel Level
}

func (p *probe) Enabled(l Level) bool {
	return l <= p.minLevel
}

func (p *probe) String() string {
	return p.text
}
