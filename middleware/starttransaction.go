package middleware

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger-apm/envelope"
	"go.elastic.co/apm"
)

// StartTransaction starts a transaction using W3C Trace-Context headers from the envelope
func StartTransaction(tracer *apm.Tracer, name string, typ string) *startTransaction {
	return &startTransaction{
		tracer: tracer,
		name:   name,
		typ:    typ,
	}
}

type startTransaction struct {
	tracer *apm.Tracer
	name   string
	typ    string
}

func (w *startTransaction) Handle(ctx context.Context, _ messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	tx := envelope.StartTransaction(w.tracer, w.name, w.typ, e)
	defer tx.End()

	ctx = apm.ContextWithTransaction(ctx, tx)

	next(ctx, e)
}
