package trace

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	Tracer trace.Tracer
	tp     *sdktrace.TracerProvider
}

type AppSpan interface {
	End()
	RecordError(err error)
	SetString(key string, value string)
	SetInt(key string, value int)
	SetFloat(key string, value float64)
	SetBool(key string, value bool)
}

type appSpan struct {
	ctx  context.Context
	span trace.Span
}

func newAppSpan(ctx context.Context, span trace.Span) AppSpan {
	return &appSpan{
		ctx:  ctx,
		span: span,
	}
}

var singleTracer *Tracer

func NewTracer(serviceName string, address string) *Tracer {
	ctx := context.Background()

	// Create http exporter
	httpExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(address), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry HTTP exporter: %v", err)
	}

	// create span processors for each exporter
	httpSpanProcessor := sdktrace.NewBatchSpanProcessor(httpExporter)

	// Create tracer provider with span processors
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(httpSpanProcessor),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	singleTracer = &Tracer{
		Tracer: tp.Tracer(serviceName),
		tp:     tp,
	}

	log.Println(fmt.Sprintf("Jaeger tracer connected on %s", address))
	return singleTracer
}

func (t *Tracer) Start(ctx context.Context, name string) (context.Context, AppSpan) {
	ctx, span := t.Tracer.Start(ctx, name)
	return ctx, newAppSpan(ctx, span)
}

func (a *appSpan) RecordError(err error) {
	a.span.RecordError(
		err, trace.WithAttributes(attribute.String("error", err.Error())),
	)
	a.span.SetStatus(codes.Error, err.Error())
}

func (a *appSpan) End() {
	a.span.End()
}

func (a *appSpan) SetString(key string, value string) {
	a.span.SetAttributes(attribute.String(key, value))
}

func (a *appSpan) SetInt(key string, value int) {
	a.span.SetAttributes(attribute.Int(key, value))
}

func (a *appSpan) SetFloat(key string, value float64) {
	a.span.SetAttributes(attribute.Float64(key, value))
}

func (a *appSpan) SetBool(key string, value bool) {
	a.span.SetAttributes(attribute.Bool(key, value))
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if err := t.tp.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not shutdown OpenTelemetry tracer: %v", err)
	}
	return nil
}

func GetTracer() *Tracer {
	if singleTracer == nil {
		singleTracer = NewTracer("test", "4813")
	}
	return singleTracer
}
