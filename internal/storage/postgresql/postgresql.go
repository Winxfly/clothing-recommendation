package postgresql

import (
	"clothing-recommendation/internal/config"
	"clothing-recommendation/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB *sql.DB
}

func New(cfgStorage config.StoragePath) (*Storage, error) {
	const op = "storage.postgresql.New"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfgStorage.Username, cfgStorage.Password, cfgStorage.Host,
		cfgStorage.Port, cfgStorage.Database, cfgStorage.SSLMode)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w:", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w:", op, err)
	}

	return &Storage{DB: db}, nil
}

func roundToNearFive(n int) int {
	return ((n + 2) / 5) * 5
}

func (s *Storage) GetRecommendation(rawTemp, rawWind int) (string, error) {
	const op = "storage.postgresql.GetRecommendation"

	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	temp := roundToNearFive(rawTemp)
	var tempRec string
	err = tx.QueryRowContext(ctx,
		`SELECT recommendation FROM temperature WHERE temperature = $1`,
		temp,
	).Scan(&tempRec)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrTemperatureNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	wind := roundToNearFive(rawWind)
	var windRec string
	err = tx.QueryRowContext(ctx,
		`SELECT recommendation FROM wind WHERE wind_speed = $1`,
		wind,
	).Scan(&windRec)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrWindNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tempRec + windRec, nil
}
