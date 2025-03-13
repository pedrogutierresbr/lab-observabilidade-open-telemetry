package service

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"log/slog"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer      = otel.Tracer("service_a")
	URLServiceB string
)

func init() {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		log.Fatal("warning: could not load .env file", "error:", err)
	}

	URLServiceB = os.Getenv("URL_SERVICE_B")
	if URLServiceB == "" {
		log.Fatal("mandatory variable SERVICE_B_URL not defined in .env")
	}
}

func GetWeatherByZipCode(ctx context.Context, zipCode string) ([]byte, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetWeatherByZipCode")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", URLServiceB+zipCode, nil)
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
