package mtl

import (
	"os"

	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func TraceInit(serviceName string) provider.OtelProvider {
	opts := []provider.Option{
		provider.WithServiceName(serviceName),
		provider.WithInsecure(),
		provider.WithEnableMetrics(false),
	}
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		opts = append(opts, provider.WithExportEndpoint(endpoint))
	}
	return provider.NewOpenTelemetryProvider(opts...)
}
