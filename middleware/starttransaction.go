package middleware

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger-apm/envelope"
	"go.elastic.co/apm"
)

// StartTransaction starts a transaction using W3C Trace-Context headers from the envelope
func StartTransaction(tracer *apm.Tracer) *startTransaction {
	return &startTransaction{
		tracer: tracer,
	}
}

type startTransaction struct {
	tracer *apm.Tracer
}

func (w *startTransaction) Handle(ctx context.Context, _ messenger.Dispatcher, e messenger.Envelope, next messenger.NextFunc) {
	tx := envelope.StartTransaction(w.tracer, e)
	defer tx.End()

	ctx = apm.ContextWithTransaction(ctx, tx)

	next(ctx, e)
}
