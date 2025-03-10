package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize OpenTelemetry
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Get a tracer
	tracer := otel.Tracer("client-tracer")

	for {
		// Create a context with a span
		ctx, span := tracer.Start(context.Background(), "client-request")

		// Create a new request with the context
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://app:8080", nil)

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			span.End()
		} else {
			fmt.Println("Response status:", resp.Status)
			resp.Body.Close()
			span.End()
		}

		time.Sleep(5 * time.Second)
	}
}

func initTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()

	// Create OTLP exporter
	conn, err := grpc.DialContext(ctx, "collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Create the tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("sample-client"),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
