package otel1

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"time"

	//"kratos/pkg/net/tracing"
	"os"
)

func NewTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		//tracing.Metadata{},
		propagation.Baggage{},
		propagation.TraceContext{},
	)
}

func NewTracerProvider(ctx context.Context, textMapPropagator propagation.TextMapPropagator) (*trace.TracerProvider, error) {
	//meta := bc.Metadata
	//traceConf := bc.Otel.Trace

	//endpoint := "https://172.20.180.115:14268/v1/traces"
	//endpoint := "https://172.20.180.115:14268/v1/traces"
	//endpoint := "172.20.180.115:4318"
	//opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	//opts = append(opts, otlptracehttp.WithInsecure())
	//client := otlptracehttp.NewClient(opts...)
	//exp, err := otlptrace.New(ctx, client)
	//if err != nil {
	//	return nil, err
	//}

	traceExporter, _ := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	tp := sdktrace.NewTracerProvider(
		//sdktrace.WithBatcher(exp, trace.WithBatchTimeout(time.Second)),
		sdktrace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(os.Getenv("APP_NAME")),
				attribute.String("environment", os.Getenv("DEPLOY_ENV")),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(textMapPropagator)
	return tp, nil
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

type OtelLogger struct {
}
