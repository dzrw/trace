package trace

import "context"

type key int

var traceKey key

func NewContext(parent context.Context, tr Trace) context.Context {
	return context.WithValue(parent, traceKey, tr)
}

func FromContext(ctx context.Context) (s Trace, ok bool) {
	v, ok := ctx.Value(traceKey).(*traceimpl)
	return v, ok
}
