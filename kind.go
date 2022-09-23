package trace

import "errors"

var ErrKind = errors.New("unsupported kind")

// Kind is the kind of an Attr's value.
type Kind int

const (
	AnyKind Kind = iota
	BoolKind
	DurationKind
	ErrorKind
	Float64Kind
	Int64Kind
	NoErrorKind
	StringKind
	TimeKind
	Uint64Kind
)

func (k Kind) String() string {
	switch k {
	case BoolKind:
		return "bool"
	case DurationKind:
		return "time.Duration"
	case Float64Kind:
		return "float64"
	case Int64Kind:
		return "int64"
	case ErrorKind, NoErrorKind:
		return "error"
	case StringKind:
		return "string"
	case TimeKind:
		return "time.Time"
	case Uint64Kind:
		return "uint64"
	case AnyKind:
		fallthrough
	default:
		return "any"
	}
}
