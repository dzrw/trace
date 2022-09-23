package trace

import (
	"strconv"
	"strings"
)

// A Level is the importance or severity of a log event. The higher the level,
// the less important or severe the event.
type Level int

// Names for common levels.
const (
	ErrorLevel             Level = 10
	AssertionViolatedLevel Level = 11
	WarnLevel              Level = 20
	InfoLevel              Level = 30
	DebugLevel             Level = 31
)

/*
String returns a name for the level. If the level has a name, then that
name in uppercase is returned. If the level is between named values, then an
integer is appended to the uppercased name. Examples:

	WarnLevel.String() => "WARN"
	(WarnLevel-2).String() => "WARN-2"
*/
func (l Level) String() string {
	switch {
	case l < ErrorLevel:
		return strconv.Itoa(int(l))
	case l == ErrorLevel:
		return "ERROR"
	case l == AssertionViolatedLevel:
		return "ASSERT"
	case l < WarnLevel:
		sb := strings.Builder{}
		sb.WriteString("ERROR-")
		sb.WriteString(strconv.Itoa(int(l)))
		return sb.String()
	case l == WarnLevel:
		return "WARN"
	case l < InfoLevel:
		sb := strings.Builder{}
		sb.WriteString("WARN-")
		sb.WriteString(strconv.Itoa(int(l)))
		return sb.String()
	case l == InfoLevel:
		return "INFO"
	case l < DebugLevel:
		sb := strings.Builder{}
		sb.WriteString("INFO-")
		sb.WriteString(strconv.Itoa(int(l)))
		return sb.String()
	case l == DebugLevel:
		return "DEBUG"
	default:
		sb := strings.Builder{}
		sb.WriteString("DEBUG-")
		sb.WriteString(strconv.Itoa(int(l)))
		return sb.String()
	}
}
