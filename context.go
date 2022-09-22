package trace

import "context"

type key int

var traceKey key

func NewContext[T ProbeC](parent context.Context, tr Trace[T]) context.Context {
	return context.WithValue(parent, traceKey, tr)
}

func FromContext[T ProbeC](ctx context.Context) (s Trace[T], ok bool) {
	v, ok := ctx.Value(traceKey).(*traceimpl[T])
	return v, ok
}
