package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"service_a/internal/handler"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	Host       string
	URL_ZIPKIN string
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("warning: could not load .env file", "error:", err)
	}

	Host = os.Getenv("URL_ZIPKIN")
	if Host == "" {
		log.Fatal("mandatory variable URL_ZIPKIN not defined in .env")
	}

	URL_ZIPKIN = "http://" + Host + ":9411/api/v2/spans"
}

func initTracer() func() {
	exporter, err := zipkin.New(
		URL_ZIPKIN,
	)

	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("_a"),
		)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}
}

func main() {
	log.Println("Starting server.")
	cleanup := initTracer()
	defer cleanup()

	http.HandleFunc("/cep", handler.HandleZipcode)
	err := http.ListenAndServe(":8080", otelhttp.NewHandler(http.DefaultServeMux, "http-server"))
	if err != nil {
		log.Fatal("error listen and serve", "error:", err)
	}
	log.Println("Started on port", "8080")
}
