package jaeger

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitJaeger initializes and registers jaeger to global TracerProvider.
//
// The output parameter `tp` is used for waiting exported trace spans to be uploaded,
// which is useful if your program is ending, and you do not want to lose recent spans.
func InitJaeger(serviceName, endpoint string) (tp *trace.TracerProvider, err error) {
	var endpointOption jaeger.EndpointOption
	endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))

	// Create the Jaeger exporter
	exp, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}

	return initOTELTracer(exp, serviceName)
}

// InitStdOutForDevelopment initializes and registers stdout to global TracerProvider.
// It is used in local development, Don't use it online.
func InitStdOutForDevelopment(serviceName, endpoint string) (tp *trace.TracerProvider, err error) {
	exp, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	return initOTELTracer(exp, serviceName)
}

func initOTELTracer(exporter trace.SpanExporter, serviceName string) (tp *trace.TracerProvider, err error) {
	tp = trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exporter),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	defaultTextMapPropagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(defaultTextMapPropagator)
	return tp, nil
}
