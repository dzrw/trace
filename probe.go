package trace

import (
	"fmt"
)

type Probe interface {
	fmt.Stringer
}

func NewProbe(text string) Probe {
	return &probe{text}
}

type probe struct {
	text string
}

func (p *probe) String() string {
	return p.text
}
