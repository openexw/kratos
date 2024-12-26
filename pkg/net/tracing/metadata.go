package tracing

import (
	"context"
	"kratos/pkg/net/metadata"
	"os"

	"go.opentelemetry.io/otel/propagation"
)

const serviceHeader = "x-md-service-name"

// Metadata is tracing metadata propagator
type Metadata struct{}

var _ propagation.TextMapPropagator = Metadata{}

//var _ctxkey ctxKey = "kratos/pkg/net/trace.trace"

// FromContext returns the trace bound to the context, if any.
//func FromContext(ctx context.Context) (t Trace, ok bool) {
//	t, ok = ctx.Value(_ctxkey).(Trace)
//	return
//}
//
//// NewContext new a trace context.
//// NOTE: This method is not thread safe.
//func NewContext(ctx context.Context, t Trace) context.Context {
//	return context.WithValue(ctx, _ctxkey, t)
//}

// Inject sets metadata key-values from ctx into the carrier.
func (b Metadata) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	metadata.FromContext(ctx)
	carrier.Set(serviceHeader, os.Getenv("APP_NAME"))
}

// Extract returns a copy of parent with the metadata from the carrier added.
func (b Metadata) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	name := carrier.Get(serviceHeader)
	if name == "" {
		return parent
	}

	if md, ok := metadata.FromServerContext(parent); ok {
		md.Set(serviceHeader, name)
		return parent
	}
	md := metadata.New(nil)
	md.Set(serviceHeader, name)
	parent = metadata.NewServerContext(parent, md)
	return parent
}

// Fields returns the keys whose values are set with Inject.
func (b Metadata) Fields() []string {
	return []string{serviceHeader}
}
