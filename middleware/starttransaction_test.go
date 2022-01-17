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
	"testing"
)

func TestStartTransaction_Handle_given_envelope_without_trace_headers_will_start_new_transaction(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "value")
	e := envelope.FromMessage("test message")
	b := &mock.Dispatcher{}

	m := StartTransaction(apmtest.NewDiscardTracer(), "test name", "test type")

	var nextCtx context.Context
	var nextE messenger.Envelope
	m.Handle(ctx, b, e, func(ctx context.Context, e messenger.Envelope) {
		nextCtx = ctx
		nextE = e
		tx := apm.TransactionFromContext(ctx)
		if !assert.NotNil(t, tx) {
			t.Fatal("expected transaction but got nil")
		}
		assert.Equal(t, "test name", tx.Name)
		assert.Equal(t, "test type", tx.Type)
	})

	tx := apm.TransactionFromContext(nextCtx)
	if !assert.NotNil(t, tx) {
		t.Fatal("expected transaction but got nil")
	}

	assert.Nil(t, tx.TransactionData)
	assert.Same(t, e, nextE)
	assert.Equal(t, "value", ctx.Value("test"))
}

func TestStartTransaction_Handle_given_envelope_with_trace_headers_will_start_new_transaction_with_parent_from_headers(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "value")

	var e messenger.Envelope = envelope.FromMessage("test message")
	e = envelope.WithHeader(e, libenvelope.TraceparentHeader, "00-8475ef5f80f95de1f3bbcc1c7c64b7c9-8475ef5f80f95de1-01")
	e = envelope.WithHeader(e, libenvelope.TracestateHeader, "es=s:1")

	b := &mock.Dispatcher{}

	m := StartTransaction(apmtest.NewDiscardTracer(), "test name", "test type")

	var nextCtx context.Context
	var nextE messenger.Envelope
	m.Handle(ctx, b, e, func(ctx context.Context, e messenger.Envelope) {
		nextCtx = ctx
		nextE = e
		tx := apm.TransactionFromContext(ctx)
		if !assert.NotNil(t, tx) {
			t.Fatal("expected transaction but got nil")
		}
		assert.Equal(t, "test name", tx.Name)
		assert.Equal(t, "test type", tx.Type)
		assert.Equal(t, "8475ef5f80f95de1f3bbcc1c7c64b7c9", tx.TraceContext().Trace.String())
		assert.Equal(t, "8475ef5f80f95de1", tx.ParentID().String())
	})

	tx := apm.TransactionFromContext(nextCtx)
	if !assert.NotNil(t, tx) {
		t.Fatal("expected transaction but got nil")
	}

	assert.Nil(t, tx.TransactionData)
	assert.Same(t, e, nextE)
	assert.Equal(t, "value", ctx.Value("test"))
}
