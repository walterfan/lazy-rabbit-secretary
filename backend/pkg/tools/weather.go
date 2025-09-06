package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// WeatherData represents weather information from Gaode API
type WeatherData struct {
	Status    string     `json:"status"`
	Count     string     `json:"count"`
	Info      string     `json:"info"`
	Infocode  string     `json:"infocode"`
	Lives     []Live     `json:"lives"`
	Forecasts []Forecast `json:"forecasts"`
}

// Live represents real-time weather data
type Live struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	Adcode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	Winddirection    string `json:"winddirection"`
	Windpower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	Reporttime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	Humidity_float   string `json:"humidity_float"`
}

// Forecast represents forecast weather data
type Forecast struct {
	City       string `json:"city"`
	Adcode     string `json:"adcode"`
	Province   string `json:"province"`
	Reporttime string `json:"reporttime"`
	Casts      []Cast `json:"casts"`
}

// Cast represents daily forecast
type Cast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	Dayweather   string `json:"dayweather"`
	Nightweather string `json:"nightweather"`
	Daytemp      string `json:"daytemp"`
	Nighttemp    string `json:"nighttemp"`
	Daywind      string `json:"daywind"`
	Nightwind    string `json:"nightwind"`
	Daypower     string `json:"daypower"`
	Nightpower   string `json:"nightpower"`
}

// getWeatherData calls Gaode Weather API
func getWeatherData(city string, ext string) (*WeatherData, error) {
	gaodeKey := os.Getenv("LBS_API_KEY")
	if gaodeKey == "" {
		return nil, fmt.Errorf("LBS_API_KEY environment variable not set")
	}

	baseURL := os.Getenv("LBS_API_URL")
	if baseURL == "" {
		baseURL = "https://restapi.amap.com/v3/weather/weatherInfo"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build request URL
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("key", gaodeKey)
	q.Add("city", city)
	q.Add("extensions", ext)
	q.Add("output", "JSON")
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON
	var weatherData WeatherData
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Check API response status
	if weatherData.Status != "1" {
		return nil, fmt.Errorf("API error: %s (code: %s)", weatherData.Info, weatherData.Infocode)
	}

	return &weatherData, nil
}

// GetWeatherFunction returns the weather function definition
func GetWeatherFunction() FunctionDefinition {
	return FunctionDefinition{
		Name:        "get_weather",
		Description: "Get weather of a city. User must supply a city first.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"city": map[string]interface{}{
					"type":        "string",
					"description": "城市代码(adcode)或中文城市名，例如：北京、上海 (City code or Chinese city name, e.g., 北京, 上海)",
				},
				"ext": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"base", "all"},
					"description": "base=实时天气，all=预报天气，默认base (base=real-time weather, all=forecast weather, default: base)",
					"default":     "base",
				},
			},
			"required": []string{"city"},
		},
		Handler: func(args map[string]interface{}) (interface{}, error) {
			city, ok := args["city"].(string)
			if !ok || city == "" {
				return nil, fmt.Errorf("city parameter is required")
			}

			ext := "base"
			if e, ok := args["ext"].(string); ok && e != "" {
				ext = e
			}

			weatherData, err := getWeatherData(city, ext)
			if err != nil {
				return nil, fmt.Errorf("failed to get weather data: %v", err)
			}

			// Format response for easier consumption
			if ext == "base" && len(weatherData.Lives) > 0 {
				live := weatherData.Lives[0]
				return map[string]interface{}{
					"city":        live.City,
					"province":    live.Province,
					"weather":     live.Weather,
					"temperature": live.Temperature,
					"humidity":    live.Humidity,
					"wind":        fmt.Sprintf("%s %s级", live.Winddirection, live.Windpower),
					"reporttime":  live.Reporttime,
					"type":        "real-time",
				}, nil
			} else if ext == "all" && len(weatherData.Forecasts) > 0 {
				forecast := weatherData.Forecasts[0]
				var casts []map[string]interface{}
				for _, cast := range forecast.Casts {
					casts = append(casts, map[string]interface{}{
						"date":         cast.Date,
						"week":         cast.Week,
						"dayweather":   cast.Dayweather,
						"nightweather": cast.Nightweather,
						"daytemp":      cast.Daytemp,
						"nighttemp":    cast.Nighttemp,
						"daywind":      cast.Daywind,
						"nightwind":    cast.Nightwind,
					})
				}
				return map[string]interface{}{
					"city":       forecast.City,
					"province":   forecast.Province,
					"reporttime": forecast.Reporttime,
					"forecasts":  casts,
					"type":       "forecast",
				}, nil
			}

			return weatherData, nil
		},
	}
}
