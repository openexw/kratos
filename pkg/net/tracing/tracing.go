package tracing

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	KratosTraceID = "kratos-trace-id"
)

// Option is tracing option.
type Option func(*options)

type options struct {
	tracerName     string
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

// WithPropagator with tracer propagator.
func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(opts *options) {
		opts.propagator = propagator
	}
}

// WithTracerProvider with tracer provider.
// By default, it uses the global provider that is set by otel.SetTracerProvider(provider).
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(opts *options) {
		opts.tracerProvider = provider
	}
}

// WithTracerName with tracer name
func WithTracerName(tracerName string) Option {
	return func(opts *options) {
		opts.tracerName = tracerName
	}
}

//
//// Server returns a new server middleware for OpenTelemetry.
//func Server(opts ...Option) middleware.Middleware {
//	tracer := NewTracer(trace.SpanKindServer, opts...)
//	return func(handler middleware.Handler) middleware.Handler {
//		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
//			if tr, ok := transport.FromServerContext(ctx); ok {
//				var span trace.Span
//				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
//				SetServerSpan(ctx, span, req)
//				defer func() { tracer.End(ctx, span, reply, err) }()
//			}
//			return handler(ctx, req)
//		}
//	}
//}
//
//// Client returns a new client middleware for OpenTelemetry.
//func Client(opts ...Option) middleware.Middleware {
//	tracer := NewTracer(trace.SpanKindClient, opts...)
//	return func(handler middleware.Handler) middleware.Handler {
//		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
//			if tr, ok := transport.FromClientContext(ctx); ok {
//				var span trace.Span
//				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
//				SetClientSpan(ctx, span, req)
//				defer func() { tracer.End(ctx, span, reply, err) }()
//			}
//			return handler(ctx, req)
//		}
//	}
//}

// TraceID returns a traceid valuer.
func TraceID(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}

// SpanID returns a spanid valuer.
func SpanID(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
		return span.SpanID().String()
	}
	return ""
}
