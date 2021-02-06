package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"
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

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	cityResults := make(chan Weather, len(cities))
	for _, city := range cities {
		city := city
		eg.Go(func() error {
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
			req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
			if err != nil {
				return fmt.Errorf("failed to create new request: %w", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("unexpected error getting weather for city: %w", err)
			}
			defer resp.Body.Close()
			weather := Weather{
				City: city,
			}
			err = json.NewDecoder(resp.Body).Decode(&weather)
			if err != nil {
				return fmt.Errorf("unexpected error parsing weather response: %w", err)
			}
			cityResults <- weather
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Println("Failed to get weather", err)
		return
	}
	close(cityResults)
	for result := range cityResults {
		log.Printf("The temperature in %s was %v degrees C", result.City, result.Main.Temp)
	}
}

// Weather describes the JSON structure returned from the weather API
type Weather struct {
	City  string
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
