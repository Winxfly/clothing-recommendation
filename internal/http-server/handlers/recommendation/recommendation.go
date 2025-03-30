package recommendation

import (
	"clothing-recommendation/internal/http-server/middleware/logger"
	"clothing-recommendation/internal/lib/api/response"
	loggerSlog "clothing-recommendation/internal/lib/logger/slog"
	"clothing-recommendation/internal/storage"
	"clothing-recommendation/internal/weather"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type Request struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Response struct {
	response.Response
	Recommendation string `json:"recommendation,omitempty"`
}

type Recommer interface {
	GetRecommendation(temperature int, wind int) (string, error)
}

func New(log *slog.Logger, getRec Recommer, weatherClient *weather.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.recommendation.New"
		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", logger.GetRequestID(ctx)),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", loggerSlog.Err(err))
			response.JSON(w, http.StatusBadRequest,
				response.Error("invalid request format"))

			return
		}

		temp, wind, err := weatherClient.GetWeather(ctx, req.Latitude, req.Longitude)
		if err != nil {
			log.Error("failed to get weather", loggerSlog.Err(err))
			response.JSON(w, http.StatusInternalServerError,
				response.Error("weather service error"))

			return
		}

		rec, err := getRec.GetRecommendation(int(temp), int(wind))
		if err != nil {
			if errors.Is(err, storage.ErrTemperatureNotFound) || errors.Is(err, storage.ErrWindNotFound) {
				response.JSON(w, http.StatusNotFound,
					response.Error("recommendation not found"))

				return
			}

			log.Error("storage error", loggerSlog.Err(err))
			response.JSON(w, http.StatusInternalServerError,
				response.Error("internal error"))

			return
		}

		rec = rec + fmt.Sprintf(" Температура:%.1f °C Ветер:%.1f м/с", temp, wind)

		response.JSON(w, http.StatusOK, Response{
			Response: response.OK(rec),
		})
	}
}
