package geocode

import (
	"clothing-recommendation/internal/http-server/middleware/logger"
	"clothing-recommendation/internal/lib/api/response"
	loggerSlog "clothing-recommendation/internal/lib/logger/slog"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type GeoClient struct {
	baseURL string
}

func New(baseURL string) *GeoClient {
	return &GeoClient{baseURL: baseURL}
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *GeoClient) Search(ctx context.Context, query string) ([]Location, error) {
	url := fmt.Sprintf("%s?name=%s&count=5&language=ru",
		c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to geocoding: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to geocoding: invalid status code: %d", resp.StatusCode)
	}

	var result struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	locations := make([]Location, 0, len(result.Results))
	for _, r := range result.Results {
		locations = append(locations, Location{
			Name:      r.Name,
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		})
	}

	return locations, nil
}

func Handler(log *slog.Logger, client *GeoClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.geocode.Handler"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", logger.GetRequestID(r.Context())),
		)

		query := r.URL.Query().Get("query")
		if query == "" {
			response.JSON(w, http.StatusBadRequest,
				response.Error("bad request"))
			return
		}

		locations, err := client.Search(r.Context(), query)
		if err != nil {
			log.Error("geocoding failed", loggerSlog.Err(err))
			response.JSON(w, http.StatusInternalServerError,
				response.Error("location search failed"))

			return
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"results": locations,
		})
	}
}
