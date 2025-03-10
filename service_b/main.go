package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	otelzipkin "go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type CEPWeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func fetchCityData(cep string) (string, error) {
	apiUrl := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", fmt.Errorf("data could not be fetched")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	city := data["localidade"].(string)
	return city, nil
}

func fetchWeatherData(city string) (CEPWeatherResponse, error) {
	apiUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", "e4445b7bd7f54bd580c222033251502", city)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return CEPWeatherResponse{}, fmt.Errorf("failed to fetch temperature: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CEPWeatherResponse{}, fmt.Errorf("failed to fetch temperature, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CEPWeatherResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return CEPWeatherResponse{}, fmt.Errorf("failed to decode JSON: %v", err)
	}

	current, ok := data["current"].(map[string]interface{})
	if !ok {
		return CEPWeatherResponse{}, fmt.Errorf("invalid response format")
	}

	tempC, ok := current["temp_c"].(float64)
	if !ok {
		return CEPWeatherResponse{}, fmt.Errorf("invalid temperature format")
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return CEPWeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

func startTracer() trace.Tracer {
	exporter, err := otelzipkin.New("http://localhost:9441/api/v2/spans")
	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Tracer("service_b")
}

func handleCEPWeatherRequest(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "handleCEPWeatherRequest")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := fetchCityData(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	resp, err := fetchWeatherData(city)
	if err != nil {
		http.Error(w, "failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	tracer = startTracer()

	http.Handle("/getCEPWeatherData", http.HandlerFunc(handleCEPWeatherRequest))
	fmt.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
