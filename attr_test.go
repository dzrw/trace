package trace_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/dzrw/trace"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	a := trace.Bool("foo", true)
	require.NotNil(t, a)
	require.True(t, a.Bool())
}

func TestFormatBool(t *testing.T) {
	a := trace.Bool("foo", true)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=true", s.String())
}

func TestDuration(t *testing.T) {
	a := trace.Duration("foo", time.Hour)
	require.NotNil(t, a)
	require.Equal(t, time.Hour, a.Duration())
}

func TestFormatDuration(t *testing.T) {
	a := trace.Duration("foo", time.Hour)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=1h0m0s", s.String())
}

func TestFloat64(t *testing.T) {
	a := trace.Float64("foo", math.Pi)
	require.NotNil(t, a)
	require.Equal(t, math.Pi, a.Float64())
}

func TestFormatFloat64(t *testing.T) {
	a := trace.Float64("foo", math.Pi)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=3.141592653589793", s.String())
}

func TestInt64(t *testing.T) {
	a := trace.Int64("foo", math.MaxInt64)
	require.NotNil(t, a)
	require.Equal(t, int64(math.MaxInt64), a.Int64())
}

func TestFormatInt64(t *testing.T) {
	a := trace.Int64("foo", math.MaxInt64)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=9223372036854775807", s.String())
}

func TestString(t *testing.T) {
	a := trace.String("foo", "bar")
	require.NotNil(t, a)
	require.Equal(t, "bar", a.String())
}

func TestFormatString(t *testing.T) {
	a := trace.String("foo", "bar")
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=bar", s.String())
}

func TestTime(t *testing.T) {
	then := time.Now()
	a := trace.Time("foo", then)
	require.NotNil(t, a)
	require.Equal(t, then, a.Time())
}

func TestFormatTime(t *testing.T) {
	value := time.Now()
	expected := fmt.Sprintf("foo=%s", value.Format(trace.RFC3339Milli))

	a := trace.Time("foo", value)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, expected, s.String())
}

func TestUint64(t *testing.T) {
	a := trace.Uint64("foo", math.MaxUint64)
	require.NotNil(t, a)
	require.Equal(t, uint64(math.MaxUint64), a.Uint64())
}

func TestFormatUint64(t *testing.T) {
	a := trace.Uint64("foo", math.MaxUint64)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, "foo=18446744073709551615", s.String())
}

type astruct struct{}

func (a *astruct) String() string {
	return "hello, world"
}

func TestAny(t *testing.T) {
	val := &astruct{}
	a := trace.Any("foo", val)
	require.NotNil(t, a)
	require.NotNil(t, a.Value())
	require.True(t, a.HasValue())

	got, ok := a.Value().(*astruct)
	require.True(t, ok)
	require.NotNil(t, got)
	require.Equal(t, val, got)
}

func TestFormatAny(t *testing.T) {
	val := &astruct{}
	expected := fmt.Sprintf("foo=%s", val)

	a := trace.Any("foo", val)
	s := &State{}
	a.Format(s, ' ')
	require.Equal(t, expected, s.String())
}
