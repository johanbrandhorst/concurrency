package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var citiesFlag = flag.String("cities", "", "comma-separated list of cities to request weather for")
var apiKey = flag.String("api-key", "", "The API key used for the API")

func main() {
	flag.Parse()
	if *citiesFlag == "" {
		log.Fatalln("must provide at least 1 city")
	}
	if *apiKey == "" {
		log.Fatalln("API key is required")
	}
	cities := strings.Split(*citiesFlag, ",")

	for _, city := range cities {
		q := url.Values{
			"q":     []string{city},
			"appid": []string{*apiKey},
			"units": []string{"metric"},
		}
		apiURL := &url.URL{
			Scheme:   "https",
			Host:     "api.openweathermap.org",
			Path:     "/data/2.5/weather",
			RawQuery: q.Encode(),
		}
		resp, err := http.Get(apiURL.String())
		if err != nil {
			log.Println("unexpected error getting weather for city:", err)
			continue
		}
		defer resp.Body.Close()
		var weather Weather
		err = json.NewDecoder(resp.Body).Decode(&weather)
		if err != nil {
			log.Println("unexpected error parsing weather response:", err)
			continue
		}
		log.Printf("The temperature in %s is %v degrees C", city, weather.Main.Temp)
	}
}

// Weather describes the JSON structure returned from the weather API
type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}
