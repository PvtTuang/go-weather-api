package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type WeatherData struct {
	CurrentCondition []struct {
		TempC       string `json:"temp_C"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		Humidity string `json:"humidity"`
	} `json:"current_condition"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := mux.Vars(r)["city"]
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", city)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data WeatherData

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding weather data", http.StatusInternalServerError)
		return
	}

	output := map[string]string{
		"city":        city,
		"temp_C":      data.CurrentCondition[0].TempC,
		"description": data.CurrentCondition[0].WeatherDesc[0].Value,
		"humidity":    data.CurrentCondition[0].Humidity,
	}

	json.NewEncoder(w).Encode(output)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{city}", weatherHandler).Methods("GET")

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
