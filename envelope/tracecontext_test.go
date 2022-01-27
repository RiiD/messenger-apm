package envelope

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger/envelope"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm"
	"go.elastic.co/apm/apmtest"
	"go.elastic.co/apm/module/apmhttp"
	"testing"
)

func TestWithTraceContext_given_context_with_transaction_it_should_set_trace_context_headers_with_its_value(t *testing.T) {
	tracer := apmtest.NewDiscardTracer()
	tx1 := tracer.StartTransaction("Test1", "test")
	tx2 := tracer.StartTransaction("Test2", "test")

	ctx1 := apm.ContextWithTransaction(context.Background(), tx1)
	ctx2 := apm.ContextWithTransaction(context.Background(), tx2)

	e := WithTraceContext(ctx1, envelope.FromMessage("test message"))
	e = WithTraceContext(ctx2, e)

	traceparent := e.Header(TraceparentHeader)
	tracestate := e.Header(TracestateHeader)

	assert.Equal(t, traceparent, []string{apmhttp.FormatTraceparentHeader(tx2.TraceContext())})
	assert.Equal(t, tracestate, []string{tx2.TraceContext().State.String()})
}

func TestWithTraceContext_given_context_without_transaction_it_should_return_same_envelope_without_change(t *testing.T) {
	ctx := context.Background()

	rootE := envelope.FromMessage("test message")
	e := WithTraceContext(ctx, rootE)

	assert.Same(t, rootE, e)
}

func TestStartTransaction_given_envelope_without_trace_context_headers_when_start_transaction_called_it_will_return_new_transaction(t *testing.T) {
	tracer := apmtest.NewDiscardTracer()
	e := envelope.FromMessage("test message")

	tx := StartTransaction(tracer, "Test transaction", "test", e)

	assert.NotEmpty(t, tx)
	assert.Equal(t, "Test transaction", tx.Name)
	assert.Equal(t, "test", tx.Type)
	assert.Equal(t, "0000000000000000", tx.ParentID().String())
}

func TestStartTransaction_given_envelope_with_trace_context_headers_when_start_transaction_called_it_will_return_transaction_with_proper_parent(t *testing.T) {
	tracer := apmtest.NewDiscardTracer()
	var e messenger.Envelope = envelope.FromMessage("test message")
	e = envelope.WithHeader(e, TraceparentHeader, "00-8475ef5f80f95de1f3bbcc1c7c64b7c9-8475ef5f80f95de1-01")
	e = envelope.WithHeader(e, TracestateHeader, "es=s:1")

	tx := StartTransaction(tracer, "Test transaction", "test", e)

	assert.NotEmpty(t, tx)
	assert.Equal(t, "Test transaction", tx.Name)
	assert.Equal(t, "test", tx.Type)
	assert.Equal(t, "8475ef5f80f95de1", tx.ParentID().String())
	assert.Equal(t, "8475ef5f80f95de1f3bbcc1c7c64b7c9", tx.TraceContext().Trace.String())
}
