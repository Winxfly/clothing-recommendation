package postgresql

import (
	"clothing-recommendation/internal/config"
	"clothing-recommendation/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath config.StoragePath) (*Storage, error) {
	const op = "storage.postgresql.New"
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		storagePath.Username, storagePath.Password, storagePath.Host,
		storagePath.Port, storagePath.Database, storagePath.SSLMode)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w:", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w:", op, err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) GetRecommendation(temperature int, wind int) (string, error) {
	const op = "storage.postgresql.GetRecommendation"

	stmt, err := s.DB.Prepare("SELECT recommendation FROM temperature WHERE temperature = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var temperatureRecommendation string
	err = stmt.QueryRow(temperature).Scan(&temperatureRecommendation)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrTemperatureNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	stmt, err = s.DB.Prepare("SELECT recommendation FROM wind WHERE wind_speed = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var windRecommendation string
	err = stmt.QueryRow(wind).Scan(&windRecommendation)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrWindNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return temperatureRecommendation + windRecommendation, nil
}
