package tracing

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"time"
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
//func WithTracerProvider(provider trace.TracerProvider) Option {
//	return func(opts *options) {
//		opts.tracerProvider = provider
//	}
//}

// WithTracerName with tracer name
func WithTracerName(tracerName string) Option {
	return func(opts *options) {
		opts.tracerName = tracerName
	}
}

func NewTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		Metadata{},
		propagation.Baggage{},
		propagation.TraceContext{},
	)
}

type ShutdownFn func(context.Context) error

type OtelOption func(*otelOptions)

type otelOptions struct {
	textMapPropagator propagation.TextMapPropagator
	environment       string
	export            sdktrace.SpanExporter
}

func WithTextMapPropagator(textMapPropagator propagation.TextMapPropagator) OtelOption {
	return func(o *otelOptions) {
		o.textMapPropagator = textMapPropagator
	}
}

func WithEnvironment(env string) OtelOption {
	return func(o *otelOptions) {
		o.environment = env
	}
}

func WithExport(ex sdktrace.SpanExporter) OtelOption {
	return func(o *otelOptions) {
		o.export = ex
	}
}

func NewHttpExporter(ctx context.Context, endpoint string, insecure bool) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)
	return otlptrace.New(ctx, client)
}

func NewStdoutExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint())
}

// Init
// endpoint := "172.20.180.115:4318"
func Init(appname string, opt ...OtelOption) (*sdktrace.TracerProvider, []ShutdownFn, error) {
	exporter, err := NewStdoutExporter()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create stdout exporter")
	}
	op := otelOptions{
		textMapPropagator: NewTextMapPropagator(),
		environment:       "prod",
		export:            exporter,
	}
	for _, o := range opt {
		o(&op)
	}
	fns := make([]ShutdownFn, 0)

	res, err := resource.New(context.Background(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appname),
			semconv.DeploymentEnvironmentKey.String(op.environment),
		))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create resource")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(op.export, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(res),
		//sdktrace.WithIDGenerator()
	)
	fns = append(fns, tp.Shutdown)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(op.textMapPropagator)
	return tp, fns, nil
}

//
//func Init() (shutdown func(context.Context) error, err error) {
//	var shutdownFuncs []func(context.Context) error
//
//	// shutdown函数调用通过shutdownFuncs注册的清理函数，并将错误连接在一起。
//	// 每个注册的清理函数将被调用一次。
//	shutdown = func(ctx context.Context) error {
//		var err error
//		for _, fn := range shutdownFuncs {
//			err = errors.WithMessage(fn(ctx), "\n")
//		}
//		shutdownFuncs = nil
//		return err
//	}
//
//	// resource
//	//res, err := resource.Merge(resource.Default(),
//	//	resource.NewWithAttributes("sdsdd",
//	//		semconv.ServiceNameKey.String("otel-demo"),
//	//		semconv.ServiceVersionKey.String("v1.0.0"),
//	//	))
//	//if err != nil {
//	//	log.Error("sds err: %v", err)
//	//	return
//	//}
//
//	// propagation
//	propagator := propagation.NewCompositeTextMapPropagator(
//		propagation.TraceContext{},
//		propagation.Baggage{},
//	)
//	otel.SetTextMapPropagator(propagator)
//
//	// trace
//	traceProvider := trace.NewTracerProvider(
//		trace.WithBatcher(traceExporter,
//			// 默认是5秒，这里设置为1秒以演示目的。
//			trace.WithBatchTimeout(time.Second)),
//		//trace.WithResource(res),
//	)
//	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
//	otel.SetTracerProvider(traceProvider)
//	return
//}
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
