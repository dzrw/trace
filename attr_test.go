package trace_test

import (
	"errors"
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
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "true", v)
}

func TestDuration(t *testing.T) {
	a := trace.Duration("foo", time.Hour)
	require.NotNil(t, a)
	require.Equal(t, time.Hour, a.Duration())
}

func TestFormatDuration(t *testing.T) {
	a := trace.Duration("foo", time.Hour)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "1h0m0s", v)
}

func TestError(t *testing.T) {
	err := errors.New("test error")
	a := trace.Error(err)
	require.NotNil(t, a)
	require.Equal(t, err, a.Error())
}

func TestFormatError(t *testing.T) {
	err := errors.New("test error")
	a := trace.Error(err)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, err.Error(), v)
}

func TestFloat64(t *testing.T) {
	a := trace.Float64("foo", math.Pi)
	require.NotNil(t, a)
	require.Equal(t, math.Pi, a.Float64())
}

func TestFormatFloat64(t *testing.T) {
	a := trace.Float64("foo", math.Pi)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "3.141592653589793", v)
}

func TestInt64(t *testing.T) {
	a := trace.Int64("foo", math.MaxInt64)
	require.NotNil(t, a)
	require.Equal(t, int64(math.MaxInt64), a.Int64())
}

func TestFormatInt64(t *testing.T) {
	a := trace.Int64("foo", math.MaxInt64)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "9223372036854775807", v)
}

func TestString(t *testing.T) {
	a := trace.String("foo", "bar")
	require.NotNil(t, a)
	require.Equal(t, "bar", a.String())
}

func TestFormatString(t *testing.T) {
	a := trace.String("foo", "bar")
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "bar", v)
}

func TestTime(t *testing.T) {
	then := time.Now()
	a := trace.Time("foo", then)
	require.NotNil(t, a)
	require.Equal(t, then, a.Time())
}

func TestFormatTime(t *testing.T) {
	value := time.Now()
	expected := value.Format(trace.RFC3339Milli)

	a := trace.Time("foo", value)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, expected, v)
}

func TestUint64(t *testing.T) {
	a := trace.Uint64("foo", math.MaxUint64)
	require.NotNil(t, a)
	require.Equal(t, uint64(math.MaxUint64), a.Uint64())
}

func TestFormatUint64(t *testing.T) {
	a := trace.Uint64("foo", math.MaxUint64)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, "18446744073709551615", v)
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
	a := trace.Any("foo", val)
	k, v := a.Format()
	require.Equal(t, a.Key(), k)
	require.Equal(t, val.String(), v)
}
