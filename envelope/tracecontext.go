package envelope

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger/envelope"
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

// WithTraceContext adds W3C Trace-Context headers from transaction in context.Context
func WithTraceContext(ctx context.Context, e messenger.Envelope) messenger.Envelope {
	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		return e
	}

	e = envelope.WithoutHeader(e, TraceparentHeader)
	e = envelope.WithoutHeader(e, TracestateHeader)

	tc := tx.TraceContext()
	traceparentValue := apmhttp.FormatTraceparentHeader(tc)
	tracestateValue := tc.State.String()

	e = envelope.WithHeader(e, TraceparentHeader, traceparentValue)
	e = envelope.WithHeader(e, TracestateHeader, tracestateValue)

	return e
}

// StartTransaction start transaction using W3C Trace-Context headers in the envelope
func StartTransaction(tracer *apm.Tracer, name string, typ string, e messenger.Envelope) *apm.Transaction {
	traceparentValue, _ := e.LastHeader(TraceparentHeader)
	traceContext, _ := apmhttp.ParseTraceparentHeader(traceparentValue)

	tracestateValue, _ := e.LastHeader(TracestateHeader)
	traceContext.State, _ = apmhttp.ParseTracestateHeader(tracestateValue)

	tx := tracer.StartTransactionOptions(name, typ, apm.TransactionOptions{TraceContext: traceContext})

	return tx
}
