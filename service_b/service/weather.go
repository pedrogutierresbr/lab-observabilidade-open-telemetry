package services

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/url"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type WeatherAPIResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	CelsiusTemperature    float64 `json:"temp_c"`
	FarhenheitTemperature float64 `json:"temp_f"`
}

var (
	weatherApiURL string
	apiKey        string
)

func init() {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		log.Fatal("warning: could not load .env file", "error:", err)
	}

	weatherApiURL = os.Getenv("URL_WEATHERAPI")
	if weatherApiURL == "" {
		log.Fatal("mandatory variable URL_WEATHERAPI not defined in .env")
	}

	apiKey = os.Getenv("WHEATER_API_KEY")
	if apiKey == "" {
		log.Fatal("mandatory variable WHEATER_API_KEY not defined in .env")
	}
}

func GetWeatherByCity(ctx context.Context, v ViaCEPResponse) (WeatherAPIResponse, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, "GetWeatherByCity")
	defer span.End()

	params := map[string]string{
		"key": apiKey,
		"q":   v.Localidade,
		"aqi": "no",
	}

	u, err := url.Parse(weatherApiURL)
	if err != nil {
		slog.Error("error parsing URL", "url", weatherApiURL, "error", err)
		return WeatherAPIResponse{}, err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		slog.Error("unable to make new request with context", "ctx", ctx, "error", err)
		return WeatherAPIResponse{}, err
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error sending request", "query", u.RawQuery, "error", err)
		return WeatherAPIResponse{}, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response", "response:", resp.Body, "error", err)
		return WeatherAPIResponse{}, err
	}

	var weather WeatherAPIResponse
	err = json.Unmarshal(result, &weather)
	if err != nil {
		slog.Error("error unmarshal result", "result:", string(result), "error", err)
		return WeatherAPIResponse{}, err
	}
	return weather, nil
}
