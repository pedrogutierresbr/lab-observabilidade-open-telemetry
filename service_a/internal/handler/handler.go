package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"service_a/internal/service"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TemperatureInputDTO struct {
	Zipcode string `json:"cep"`
}

type TemperatureOutputDTO struct {
	City       string `json:"city"`
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

var tracer = otel.Tracer("service_a")

func HandleZipcode(w http.ResponseWriter, r *http.Request) {
	var span trace.Span
	ctx, span := tracer.Start(r.Context(), "handleZipcode")
	defer span.End()

	var inputDto TemperatureInputDTO
	err := json.NewDecoder(r.Body).Decode(&inputDto)
	if err != nil {
		slog.Error("unable to decode", "zipcode", inputDto.Zipcode, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !service.ValidZipcode(inputDto.Zipcode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	response, err := service.GetWeatherByZipCode(ctx, inputDto.Zipcode)
	if err != nil {
		http.Error(w, `"unable to fetch temperature by zipcode"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
