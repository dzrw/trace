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

type ProbeC interface {
	~int16
	Probe
}
