package apm

import (
	"context"
	"github.com/riid/messenger"
	"go.elastic.co/apm"
)

func WithSpan(name, typ string, middleware messenger.Middleware) *withSpan {
	return &withSpan{
		name:       name,
		typ:        typ,
		middleware: middleware,
	}
}

type withSpan struct {
	name       string
	typ        string
	middleware messenger.Middleware
}

func (w *withSpan) Handle(ctx context.Context, bus messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	span, spanCtx := apm.StartSpan(ctx, w.name, w.typ)
	defer span.End()

	w.middleware.Handle(spanCtx, bus, e, next)
}
