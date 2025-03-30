package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type WeatherResponse struct {
	Daily struct {
		Temperature2mMin []float64 `json:"temperature_2m_min"`
		WindSpeed10mMax  []float64 `json:"wind_speed_10m_max"`
	} `json:"daily"`
}

func (c *Client) GetWeather(ctx context.Context, lat, lon float64) (float64, float64, error) {
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&daily=temperature_2m_min,wind_speed_10m_max&forecast_days=1&wind_speed_unit=ms",
		c.baseURL, lat, lon)

	if lat == 0 || lon == 0 {
		return 0, 0, errors.New("latitude or longitude is zero")
	}

	log.Printf("Запрос для координат: lat=%f, lon=%f", lat, lon)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.Daily.Temperature2mMin) == 0 || len(data.Daily.WindSpeed10mMax) == 0 {
		return 0, 0, errors.New("no weather data available")
	}

	log.Printf("[DEBUG] Weather response: %+v", data)
	log.Printf("Данные погоды: temp=%f, wind=%f", data.Daily.Temperature2mMin[0],
		data.Daily.WindSpeed10mMax[0])
	return data.Daily.Temperature2mMin[0], data.Daily.WindSpeed10mMax[0], nil
}
