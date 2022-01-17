package middleware

import (
	"context"
	"github.com/riid/messenger"
	"go.elastic.co/apm"
)

// WithSpan wraps context passed to the next middlewares in a span
func WithSpan(name, typ string) *withSpan {
	return &withSpan{
		name: name,
		typ:  typ,
	}
}

type withSpan struct {
	name string
	typ  string
}

func (w *withSpan) Handle(ctx context.Context, _ messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	s, ctx := apm.StartSpan(ctx, w.name, w.typ)
	defer s.End()

	next(apm.ContextWithSpan(ctx, s), e)
}
