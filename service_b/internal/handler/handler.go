package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	services "service_b/service"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TemperatureOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

var (
	tracer = otel.Tracer("service_b")
)

func HandleTemperature(w http.ResponseWriter, r *http.Request) {
	var span trace.Span
	ctx, span := tracer.Start(r.Context(), "HandleTemperature")
	defer span.End()
	zipCode := strings.TrimPrefix(r.URL.Path, "/cep/")
	location, err := services.GetLocationByCEP(ctx, zipCode)
	if err != nil {
		slog.Error("failed to fetch location by zipCode", "input:", zipCode, "error", err)
		http.Error(w, "cannot find zipCode", http.StatusNotFound)
		return
	}
	temperature, err := services.GetWeatherByCity(ctx, location)
	if err != nil {
		slog.Error("failed to fetch temperature by location", "input:", location.Localidade, "error", err)
		http.Error(w, "could not get weather", http.StatusInternalServerError)
		return
	}
	formatCelcius := fmt.Sprintf("%.1f", temperature.Current.CelsiusTemperature)
	var dto TemperatureOutputDTO = TemperatureOutputDTO{
		City:       location.Localidade,
		Celsius:    formatCelcius,
		Fahrenheit: ConvertCelsiusToFahrenheit(temperature.Current.CelsiusTemperature),
		Kelvin:     ConvertCelsiusToKelvin(temperature.Current.CelsiusTemperature),
	}
	byteJson, err := json.Marshal(dto)
	if err != nil {
		slog.Error("failed to marshal dto", "dto", fmt.Sprintf("%+v", dto), "error", err)
		http.Error(w, "could not get temperature", http.StatusInternalServerError)
		return
	}
	w.Write(byteJson)
}

func ConvertCelsiusToFahrenheit(celsius float64) string {
	var f = celsius*1.8 + 32
	return fmt.Sprintf("%.1f", f)
}

func ConvertCelsiusToKelvin(celsius float64) string {
	var k = celsius + 273
	return fmt.Sprintf("%.1f", k)
}
