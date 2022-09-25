package logrus_test

import (
	"testing"

	"github.com/dzrw/trace"
	"github.com/dzrw/trace/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var SiteTestHandler = trace.Site()

func TestTextHandler(t *testing.T) {
	reg := trace.NewRegistry()
	reg.Define(SiteTestHandler, "logrus_test.SiteTestHandler")

	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	h := logrus.New(logger, false, false, reg)
	require.NotNil(t, h)
	// require.Equal(t, trace.FlagSourceInfo, h.Flags()&trace.FlagSourceInfo)
	// require.Equal(t, trace.FlagGoroutineID, h.Flags()&trace.FlagGoroutineID)

	SiteTestHandler.Install(h)
	installed, instok := SiteTestHandler.Handler()
	require.True(t, instok)
	require.Equal(t, h, installed)

	tr := SiteTestHandler.Trace(trace.Bool("test", true))
	require.NotNil(t, tr)

	tr.Debug("hello, world")
	tr.Close()

	SiteTestHandler.Gauge(99)

	require.Fail(t, "to see stdout")
}
