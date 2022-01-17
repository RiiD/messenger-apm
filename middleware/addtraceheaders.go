package middleware

import (
	"context"
	"github.com/riid/messenger"
	"github.com/riid/messenger-apm/envelope"
	"github.com/riid/messenger/middleware"
)

// AddTraceContext adds W3C Trace-Context headers to the envelope
func AddTraceContext() messenger.Middleware {
	return middleware.HandleFunc(func(ctx context.Context, _ messenger.Dispatcher, e messenger.Envelope) messenger.Envelope {
		return envelope.WithTraceContext(ctx, e)
	})
}
