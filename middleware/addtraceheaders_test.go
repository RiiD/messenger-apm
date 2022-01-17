package middleware

import (
	"context"
	"github.com/riid/messenger"
	libenvelope "github.com/riid/messenger-apm/envelope"
	"github.com/riid/messenger/envelope"
	"github.com/riid/messenger/mock"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm"
	"go.elastic.co/apm/apmtest"
	"go.elastic.co/apm/module/apmhttp"
	"testing"
)

func TestAddTraceHeaders_Handle_given_context_without_transaction_will_return_same_envelope_without_changes(t *testing.T) {
	m := AddTraceContext()
	rootE := envelope.FromMessage("test message")
	b := &mock.Dispatcher{}

	rootCtx := context.Background()

	var nextCtx context.Context
	var nextE messenger.Envelope
	m.Handle(rootCtx, b, rootE, func(ctx context.Context, e messenger.Envelope) {
		nextCtx = ctx
		nextE = e
	})

	assert.Same(t, rootCtx, nextCtx)
	assert.Same(t, rootE, nextE)
}

func TestAddTraceHeaders_Handle_given_context_with_transaction_will_add_trace_context_header_with_its_values(t *testing.T) {
	m := AddTraceContext()
	rootE := envelope.FromMessage("test message")
	b := &mock.Dispatcher{}

	tracer := apmtest.NewDiscardTracer()
	tx := tracer.StartTransaction("test transaction", "test")

	rootCtx := apm.ContextWithTransaction(context.Background(), tx)

	var nextCtx context.Context
	var nextE messenger.Envelope
	m.Handle(rootCtx, b, rootE, func(ctx context.Context, e messenger.Envelope) {
		nextCtx = ctx
		nextE = e
	})

	assert.Same(t, rootCtx, nextCtx)
	assert.True(t, nextE.Is(rootE))

	traceparent, _ := nextE.LastHeader(libenvelope.TraceparentHeader)
	tracestate, _ := nextE.LastHeader(libenvelope.TracestateHeader)

	assert.Equal(t, traceparent, apmhttp.FormatTraceparentHeader(tx.TraceContext()))
	assert.Equal(t, tracestate, tx.TraceContext().State.String())
}
