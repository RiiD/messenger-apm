package apm

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger/envelope"
	"github.com/riid/messenger/middleware"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

const (
	// TraceparentHeader is the standard W3C Trace-Context HTTP
	// header for trace propagation.
	TraceparentHeader = "Traceparent"

	// TracestateHeader is the standard W3C Trace-Context HTTP header
	// for vendor-specific trace propagation.
	TracestateHeader = "Tracestate"
)

// AddTraceContextHeaders adds W3C Trace-Context headers to the envelope
func AddTraceContextHeaders() messenger.Middleware {
	return middleware.HandleFunc(func(ctx context.Context, b messenger.Dispatcher, e messenger.Envelope) messenger.Envelope {
		tx := apm.TransactionFromContext(ctx)
		tc := tx.TraceContext()
		traceparentValue := apmhttp.FormatTraceparentHeader(tc)
		tracestateValue := tc.State.String()
		e = envelope.WithHeader(e, TraceparentHeader, traceparentValue)
		e = envelope.WithHeader(e, TracestateHeader, tracestateValue)

		return e
	})
}

// StartTransaction starts a transaction using W3C Trace-Context headers in the envelope
func StartTransaction(ctx context.Context, tracer *apm.Tracer, e messenger.Envelope, name string) context.Context {
	traceparentValue, _ := e.LastHeader(TraceparentHeader)
	traceContext, _ := apmhttp.ParseTraceparentHeader(traceparentValue)

	tracestateValue, _ := e.LastHeader(TracestateHeader)
	traceContext.State, _ = apmhttp.ParseTracestateHeader(tracestateValue)

	tx := tracer.StartTransactionOptions(name, "message", apm.TransactionOptions{TraceContext: traceContext})
	ctx = apm.ContextWithTransaction(ctx, tx)

	return ctx
}
