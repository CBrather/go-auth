package telemetry

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

func InitTracer() func(context.Context) error {
	securityOption := getGrpcSecurityOption()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			securityOption,
			otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		),
	)

	if err != nil {
		zap.L().Fatal("Tracing :: Failed setting up the OTLP trace exporter", zap.Error(err))
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "go-auth"),
		),
	)

	if err != nil {
		zap.L().Error("Tracing :: Could not set resources", zap.Error(err))
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func getGrpcSecurityOption() otlptracegrpc.Option {
	exportInsecure := os.Getenv("OTEL_EXPORTER_INSECURE_MODE")

	if strings.ToLower(exportInsecure) == "true" {
		return otlptracegrpc.WithInsecure()
	}

	return otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
}
