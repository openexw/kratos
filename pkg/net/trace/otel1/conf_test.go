package otel1

import (
	"context"
	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	os.Setenv("APP_NAME", "test-tracer")
	os.Setenv("DEPLOY_ENV", "prod")
	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	_, err := NewTracerProvider(context.Background(), propagation.NewCompositeTextMapPropagator(
		//tracing.Metadata{},
		propagation.Baggage{},
		propagation.TraceContext{},
	))
	if err != nil {
		t.Errorf("init err: %v", err)
		return
	}

	tracer := otel.Tracer("test-tracer")
	ctx, span := tracer.Start(context.Background(), "test")
	//defer span.End()
	attrs := make([]attribute.KeyValue, 0)
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		attrs = append(attrs, attribute.String("trace_id", span.TraceID().String()))
	}

	attrs = append(attrs, attribute.String("key", "value"), attribute.String("caller", "test"))
	span.SetAttributes(attrs...)
	span.End()

	var span1 trace.Span
	ctx, span1 = tracer.Start(ctx, "test1")
	//defer span.End()
	attrs1 := make([]attribute.KeyValue, 0)
	if spanx := trace.SpanContextFromContext(ctx); spanx.HasTraceID() {
		attrs1 = append(attrs1, attribute.String("trace_id", spanx.TraceID().String()))
	}

	attrs1 = append(attrs1, attribute.String("key", "value"), attribute.String("caller", "test1"))
	span1.SetAttributes(attrs1...)
	span1.End()
	time.Sleep(4 * time.Second)
}
