package trace_test

import (
	"testing"

	"github.com/dzrw/trace"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	a := trace.Bool("foo", true)
	require.NotNil(t, a)

	s := &State{}
	a.Format(s, ' ')
	require.Equal("foo=true", s.String())
}
