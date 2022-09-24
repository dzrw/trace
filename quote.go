package trace

import "strings"

func Quote(s string) string {
	if len(s) > 80 || strings.ContainsAny(s, " \t") {
		return strings.Join([]string{"\"", s, "\""}, "")
	}
	return s
}
