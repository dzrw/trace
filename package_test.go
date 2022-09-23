package trace_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dzrw/trace"
	"github.com/stretchr/testify/require"
)

type TestHandler struct {
	t *testing.T
	w *strings.Builder
	h *trace.TextHandler
}

func NewTestHandler(t *testing.T, l trace.Level) trace.Handler {
	w := strings.Builder{}
	h := trace.NewTextHandler(&w, l)
	return &TestHandler{
		t: t,
		w: &w,
		h: h,
	}
}

func (h *TestHandler) Enabled(l trace.Level) bool {
	return h.h.Enabled(l)
}

func (h *TestHandler) Log(tr trace.Trace, e *trace.EventLog) error {
	return h.h.Log(tr, e)
}

func (h *TestHandler) Count(tr trace.Trace, p trace.Probe, delta int64, attrs ...trace.Attr) error {
	return h.h.Count(tr, p, delta, attrs...)
}

func (h *TestHandler) Gauge(tr trace.Trace, p trace.Probe, value int64, attrs ...trace.Attr) error {
	return h.h.Gauge(tr, p, value, attrs...)
}

func (h *TestHandler) Flush(w io.Writer) {
	_, err := fmt.Fprintln(w, h.w.String())
	require.NoError(h.t, err)
}

func TestPackageRouterConfig(t *testing.T) {
	pkg := trace.NewPackage()
	require.NotNil(t, pkg)
	require.PanicsWithError(t, trace.ErrMustUsePackage.Error(), func() {
		pkg.Handler()
	})
	require.False(t, pkg.CaptureSourceInfo())

	r := trace.NewRouter(NewTestHandler(t, trace.DebugLevel), true)
	require.NotNil(t, r)
	require.NotNil(t, r.Handler())
	require.True(t, r.SourceInfo())

	r.Use(pkg)
	require.NotNil(t, pkg.Handler())
	require.True(t, pkg.CaptureSourceInfo())
}

var (
	probeFoo = trace.NewProbe("foo")
	probeBar = trace.NewProbe("bar")
	probeErr = trace.NewProbe("err")
)

func TestPackageUsage(t *testing.T) {
	pkg := trace.NewPackage()
	r := trace.NewRouter(NewTestHandler(t, trace.DebugLevel), false)
	r.Use(pkg)

	tr := trace.New(pkg).Title("usage test")
	tr.Debug(probeFoo, trace.Int("id", 1))

	tr.Debug(probeBar) // this should not appear
	tr.Info(probeBar)  // this should appear

	for i := 0; i < 20; i++ {
		tr.Count(probeErr, 1) // should appear after finish
		tr.Info(probeBar, trace.Int("i", i))
	}
	tr.Gauge(probeFoo, 123) // should appear after finish

	tr.Finish()

	r.Handler().(*TestHandler).Flush(os.Stdout)

	require.Fail(t, "to see stdout")
}
