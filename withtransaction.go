package apm

import (
	"context"
	"github.com/riid/messenger"
	"go.elastic.co/apm"
)

func WithTransaction(name, typ string, tracer *apm.Tracer, middleware messenger.Middleware) *withTransaction {
	return &withTransaction{
		name:       name,
		typ:        typ,
		middleware: middleware,
		tracer:     tracer,
	}
}

type withTransaction struct {
	name       string
	typ        string
	middleware messenger.Middleware
	tracer     *apm.Tracer
}

func (w *withTransaction) Handle(ctx context.Context, bus messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	parentTx := apm.TransactionFromContext(ctx)
	tx := w.tracer.StartTransactionOptions(w.name, w.typ, apm.TransactionOptions{
		TraceContext: parentTx.TraceContext(),
	})
	defer tx.End()

	headers := e.Headers()
	tx.Context.SetCustom("headers", headers)

	ctxWithTx := apm.ContextWithTransaction(ctx, tx)

	w.middleware.Handle(ctxWithTx, bus, e, next)
}
