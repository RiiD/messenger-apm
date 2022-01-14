package envelope

import (
	"context"
	"fmt"
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
	tc := tx.TraceContext()
	traceparentValue := apmhttp.FormatTraceparentHeader(tc)
	tracestateValue := tc.State.String()
	e = envelope.WithHeader(e, TraceparentHeader, traceparentValue)
	e = envelope.WithHeader(e, TracestateHeader, tracestateValue)

	return e
}

// StartTransaction start transaction using W3C Trace-Context headers in the envelope
func StartTransaction(tracer *apm.Tracer, e messenger.Envelope) *apm.Transaction {
	messageType := envelope.MessageType(e)
	if len(messageType) == 0 {
		messageType = "unknown"
	}

	traceparentValue, _ := e.LastHeader(TraceparentHeader)
	traceContext, _ := apmhttp.ParseTraceparentHeader(traceparentValue)

	tracestateValue, _ := e.LastHeader(TracestateHeader)
	traceContext.State, _ = apmhttp.ParseTracestateHeader(tracestateValue)

	name := fmt.Sprintf("%s message received", messageType)

	tx := tracer.StartTransactionOptions(name, "message", apm.TransactionOptions{TraceContext: traceContext})

	return tx
}
