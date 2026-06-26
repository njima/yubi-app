package ddtrace

import (
	"context"

	"github.com/uptrace/bun"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// BunHook implements bun's QueryHook interface and creates a Datadog APM tracing span
// around each query execution.
type BunHook struct {
	serviceName string
}

// NewBunHook returns a BunHook with the given Datadog service name.
// serviceName is used as the service identifier in Datadog APM.
func NewBunHook(serviceName string) *BunHook {
	return &BunHook{serviceName: serviceName}
}

// BeforeQuery is called immediately before a bun query executes and starts a Datadog tracing span.
// The span records the query string as ResourceName and "postgresql" as the DB type.
// The started span is embedded in the returned context and retrieved by AfterQuery.
func (h *BunHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	_, ctx = tracer.StartSpanFromContext(ctx, "bun.query",
		tracer.ServiceName(h.serviceName),
		tracer.ResourceName(event.Query),
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.Tag(ext.DBType, "postgresql"),
	)
	return ctx
}

// AfterQuery is called after a bun query completes and finishes the tracing span.
// If event.Err is non-nil, the error is recorded on the span.
// If no span is found in the context, the function is a no-op.
func (h *BunHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return
	}
	span.Finish(tracer.WithError(event.Err))
}
