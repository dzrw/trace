package trace

import (
	"encoding"
	"fmt"
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

// HACK: This is a dirty copy of the `slog` proposal.

// An Attr is a key-value pair. It can represent some small values without an
// allocation. The zero Attr has a key of "" and a value of nil.
type Attr struct {
	key  string
	val  any
	kind Kind
}

/*
Any returns an Attr for the supplied value.

Any does not preserve the exact type of integral values. All signed integers
are converted to int64 and all unsigned integers to uint64. Similarly,
float32s are converted to float64.

However, named types are preserved. So given

	type Int int

the expression

	log.Any("k", Int(1)).Value()

will return Int(1).
*/
func Any(key string, value any) Attr {
	u := value
	return Attr{key: key, val: u, kind: AnyKind}
}

// Bool returns an Attr for a bool.
func Bool(key string, value bool) Attr {
	u := value
	return Attr{key: key, val: &u, kind: BoolKind}
}

// Duration returns an Attr for a time.Duration.
func Duration(key string, value time.Duration) Attr {
	u := value
	return Attr{key: key, val: &u, kind: DurationKind}
}

// Error returns an Attr named "error" for an error.
func Error(err error) Attr {
	return Any("error", err)
}

// Float64 returns an Attr for a floating-point number.
func Float64(key string, value float64) Attr {
	u := value
	return Attr{key: key, val: &u, kind: Float64Kind}
}

// Int converts an int to an int64 and returns an Attr with that value.
func Int(key string, value int) Attr {
	u := int64(value)
	return Attr{key: key, val: &u, kind: Int64Kind}
}

// Int64 returns an Attr for an int64.
func Int64(key string, value int64) Attr {
	u := value
	return Attr{key: key, val: &u, kind: Int64Kind}
}

// String returns a new Attr for a string.
func String(key, value string) Attr {
	return Attr{key: key, val: value, kind: StringKind}
}

// Time returns an Attr for a time.Time.
func Time(key string, value time.Time) Attr {
	u := value
	return Attr{key: key, val: &u, kind: TimeKind}
}

// Uint64 returns an Attr for a uint64.
func Uint64(key string, value uint64) Attr {
	u := value
	return Attr{key: key, val: &u, kind: Uint64Kind}
}

// Bool returns the Attr's value as a bool. It panics if the value is not a
// bool.
func (a Attr) Bool() bool {
	return *(a.val.(*bool))
}

// Duration returns the Attr's value as a time.Duration. It panics if the value
// is not a time.Duration.
func (a Attr) Duration() time.Duration {
	return *(a.val.(*time.Duration))
}

// Error returns the Attr's value as a error. It panics if the value is not
// an error.
func (a Attr) Error() error {
	return a.Value().(error)
}

// Float64 returns the Attr's value as a float64. It panics if the value is not
// a float64.
func (a Attr) Float64() float64 {
	return *(a.val.(*float64))
}

// Format returns the Attr's key and value properties as strings.
func (a Attr) Format() (key, value string) {
	key = a.key

	switch a.Kind() {
	case BoolKind:
		value = fmt.Sprint(a.Bool())
	case DurationKind:
		value = fmt.Sprint(a.Duration())
	case Float64Kind:
		value = fmt.Sprint(a.Float64())
	case Int64Kind:
		value = fmt.Sprint(a.Int64())
	case StringKind:
		value = a.String()
	case TimeKind:
		value = a.Time().Format(RFC3339Milli)
	case Uint64Kind:
		value = fmt.Sprint(a.Uint64())
	case AnyKind:
		fallthrough
	default:
		if v := a.Value(); v != nil {
			if m, ok := v.(encoding.TextMarshaler); ok {
				text, err := m.MarshalText()
				if err != nil {
					panic(err)
				}
				value = string(text)
				return
			}
			value = fmt.Sprint(v)
			return
		}
		value = ""
	}
	return
}

// HasValue returns true if the Attr has a value.
func (a Attr) HasValue() bool {
	return a.val != nil
}

// Int64 returns the Attr's value as an int64. It panics if the value is not a
// signed integer.
func (a Attr) Int64() int64 {
	return *(a.val.(*int64))
}

// Key returns the Attr's key.
func (a Attr) Key() string {
	return a.key
}

// Kind returns the Attr's Kind.
func (a Attr) Kind() Kind {
	return a.kind
}

// String returns Attr's value as a string, formatted like fmt.Sprint.
// Unlike the methods Int64, Float64, and so on, which panic if the Attr is of
// the wrong kind, String never panics.
func (a Attr) String() string {
	str, _ := a.val.(string)
	return str
}

// Time returns the Attr's value as a time.Time. It panics if the value is not
// a time.Time.
func (a Attr) Time() time.Time {
	return *(a.val.(*time.Time))
}

// Uint64 returns the Attr's value as a uint64. It panics if the value is not
// an unsigned integer.
func (a Attr) Uint64() uint64 {
	return *(a.val.(*uint64))
}

// Value returns the Attr's value as an any. If the Attr does not have a value,
// it returns nil.
func (a Attr) Value() any {
	return a.val
}

// WithKey returns an attr with the given key and the receiver's value.
func (a Attr) WithKey(key string) Attr {
	return Attr{key: key, val: a.val}
}
