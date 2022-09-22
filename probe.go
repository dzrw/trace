package trace

import (
	"errors"
	"fmt"
)

var ErrUnresolved = errUnresolved
var errUnresolved = errors.New("missing or unresolvable probe")

type Probe interface {
	fmt.Stringer
	Level() Level
}

type ProbeC interface {
	~int16
	Probe
}
