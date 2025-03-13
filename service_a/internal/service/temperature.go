package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"regexp"

	"log/slog"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer = otel.Tracer("cep-service")
	URL    string
)

func init() {
	Host := os.Getenv("URL_TEMP")
	if Host == "" {
		Host = "localhost"
	}
	URL = "http://" + Host + ":8081/cep/"
}

func GetWeatherByZipCode(ctx context.Context, zipCode string) ([]byte, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetWeatherByZipCode")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", URL+zipCode, nil)
	if err != nil {
		slog.Error("unable to make new request with context", "ctx", ctx, "error", err)
		return nil, err
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("unable to do request", "req:", req.URL.Path, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func ValidZipcode(zipCode string) bool {
	resp := regexp.MustCompile(`^\d{8}$`)
	return resp.MatchString(zipCode)
}
