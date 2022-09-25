package trace_test

import (
	"os"
	"testing"

	"github.com/dzrw/trace"
	"github.com/stretchr/testify/require"
)

var SiteTestTextHandler = trace.Site()

func TestTextHandler(t *testing.T) {
	reg := trace.NewRegistry()
	reg.Define(SiteTestTextHandler, "trace_test.SiteTestTextHandler")

	h := trace.NewTextHandler(os.Stdout, trace.DebugLevel, false, false, reg)
	require.NotNil(t, h)
	// require.Equal(t, trace.FlagSourceInfo, h.Flags()&trace.FlagSourceInfo)
	// require.Equal(t, trace.FlagGoroutineID, h.Flags()&trace.FlagGoroutineID)

	SiteTestTextHandler.Install(h)
	installed, instok := SiteTestTextHandler.Handler()
	require.True(t, instok)
	require.Equal(t, h, installed)

	tr := SiteTestTextHandler.Trace(trace.Bool("test", true))
	require.NotNil(t, tr)

	tr.Debug("hello, world")
	tr.Close()

	SiteTestTextHandler.Gauge(99)

	require.Fail(t, "to see stdout")
}
