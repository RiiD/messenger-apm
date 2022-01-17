package middleware

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger/envelope"
	"github.com/riid/messenger/mock"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm"
	"go.elastic.co/apm/apmtest"
	"testing"
)

func TestSpan_Handle_given_context_with_transaction_will_call_next_with_span_in_context(t *testing.T) {
	b := &mock.Dispatcher{}

	tracer := apmtest.NewDiscardTracer()
	tx := tracer.StartTransaction("test", "test")
	rootCtx := apm.ContextWithTransaction(context.Background(), tx)

	e := envelope.FromMessage("test message")

	var nextE messenger.Envelope
	var span *apm.Span

	m := WithSpan("test name", "test type")
	m.Handle(rootCtx, b, e, func(ctx context.Context, e messenger.Envelope) {
		nextE = e
		span = apm.SpanFromContext(ctx)

		assert.Equal(t, "test name", span.Name)
		assert.Equal(t, "test type", span.Type)
		assert.Equal(t, tx.TraceContext().Span, span.ParentID())
	})

	assert.Same(t, e, nextE)
	assert.Nil(t, span.SpanData)
}
