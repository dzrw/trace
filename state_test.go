package trace_test

import (
	"bytes"
	"fmt"
)

type State struct {
	buf bytes.Buffer
}

var _ = fmt.State(&State{})

func (state *State) Write(b []byte) (n int, err error) {
	state.buf.Grow(len(b))
	return state.buf.Write(b)
}

func (state *State) Width() (wid int, ok bool) {
	return 0, false
}

func (state *State) Precision() (prec int, ok bool) {
	return 0, false
}

func (state *State) Flag(c int) bool {
	return false
}

func (state *State) Bytes() []byte {
	return state.buf.Bytes()
}

func (state *State) String() string {
	return string(state.Bytes())
}
