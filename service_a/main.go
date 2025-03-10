package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ZipCodeRequest struct {
	CEP string `json:"cep"`
}

func validateZipCode(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	if _, err := strconv.Atoi(cep); err != nil {
		return false
	}

	return true
}

func fetchCEPWeatherData(cep string) ([]byte, error) {
	apiUrl := fmt.Sprintf("http://localhost:8081/getCEPWeatherData?cep=%s", cep)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("data could not be fetched")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response, status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func handleZipCodeRequest(w http.ResponseWriter, r *http.Request) {
	var req ZipCodeRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &req)
	if err != nil || !validateZipCode(req.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	resp, err := fetchCEPWeatherData(req.CEP)
	if err != nil {
		http.Error(w, "failure to communicate with service B", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	http.Handle("/getData", http.HandlerFunc(handleZipCodeRequest))
	fmt.Println("Service A is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
