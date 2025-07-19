package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type EnvConfig struct {
	APIKey string
	BaseURL string
	Port string
}

func loadEnvConfig() *EnvConfig {
	return &EnvConfig {
		APIKey: getEnvOrDefault("API_KEY", ""),
		BaseURL: getEnvOrDefault("BASE_URL", "https://api.openweathermap.org/data/2.5/weather"),
		Port: getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var envConfig *EnvConfig

func main() {
	// Load env config after .env is loaded
	envConfig = loadEnvConfig()
	
	if envConfig.APIKey == "" {
		log.Fatal("Please set API_KEY in .env file")
	}

	// Test to see if .env is loaded correctly
	// log.Printf("API Key: %s", envConfig.APIKey)
	// log.Printf("Base URL: %s", envConfig.BaseURL)
	// log.Printf("Port: %s", envConfig.Port)

	cities := []string{"London", "Vancouver", "Tokyo", "Berlin", "Sydney", "Floridablanca", "Cebu", "Davao"}
	ch := make(chan string)
	var wg sync.WaitGroup

	startNow := time.Now()
	for _, city := range cities {
		wg.Add(1) // Increment the WaitGroup counter
		go fetchWeather(city, ch, &wg) // Fetch weather data concurrently
	}

	// Wait for all goroutines to finish and then close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
	fmt.Printf("Time taken to fetch weather data: %v\n", time.Since(startNow))
}

func fetchWeather(city string, ch chan<- string, wg *sync.WaitGroup) interface{} {
	// Ensure the WaitGroup counter is decremented when the function completes
	defer wg.Done()

	var data struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	url := fmt.Sprintf("%s?q=%s&appid=%s", envConfig.BaseURL, city, envConfig.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching weather for %s: %s\n", city, err)
		return data
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- fmt.Sprintf("Error decoding weather data for %s: %s\n", city, err)
		return data
	}

	ch <- fmt.Sprintf("This is the city %s and the temperature is %.2f", city, data.Main.Temp)
	return data
}