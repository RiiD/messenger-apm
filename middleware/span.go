package middleware

import (
	"context"
	"github.com/riid/messenger"
	"go.elastic.co/apm"
)

// Span wraps next middlewares in a span
func Span(name, typ string) *span {
	return &span{
		name: name,
		typ:  typ,
	}
}

type span struct {
	name string
	typ  string
}

func (w *span) Handle(ctx context.Context, _ messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	span, spanCtx := apm.StartSpan(ctx, w.name, w.typ)
	defer span.End()

	next(spanCtx, e)
}
