package otel

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"kratos/pkg/net/trace"
)

type Reporter struct {
	tp *tracesdk.TracerProvider
}

func (r *Reporter) WriteSpan(raw *trace.Span) error {
	//ctx := raw.Context()
	//traceID := trace.TraceID{Low: ctx.TraceID}
	//spanID := SpanID(ctx.SpanID)
	//parentID := SpanID(ctx.ParentID)
	//tags := raw.Tags()
	////log.Info("[info] write span")
	//span := &trace.Span{
	//	context:       NewSpanContext(traceID, spanID, parentID, true, nil),
	//	operationName: raw.OperationName(),
	//	startTime:     raw.StartTime(),
	//	duration:      raw.Duration(),
	//}
	//
	//span.serviceName = raw.ServiceName()
	//
	//for _, t := range tags {
	//	span.SetTag(t.Key, t.Value)
	//}
	//
	//tr := r.tp.Tracer(raw.OperationName())
	//
	//ctx, span := tr.Start(ctx, raw.ServiceName())
	//defer span.End()

	//bar(ctx)
	return nil
}

func bar(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()
}

func (r *Reporter) Close() error {
	//TODO implement me
	panic("implement me")
}

//tr := tp.Tracer("component-main")
//
//ctx, span := tr.Start(ctx, "foo")
//defer span.End()
//
//bar(ctx)
