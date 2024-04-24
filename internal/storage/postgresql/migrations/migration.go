package main

import (
	"clothing-recommendation/internal/config"
	"clothing-recommendation/internal/storage/postgresql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Data struct {
	Temperature []struct {
		Temperature    int    `json:"temperature"`
		Recommendation string `json:"recommendation"`
	} `json:"temperature"`
	Wind []struct {
		WindSpeed      int    `json:"wind_speed"`
		Recommendation string `json:"recommendation"`
	} `json:"wind"`
}

func main() {
	cfg := config.MustLoad()

	storage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Fatal("failed to initialize storage for run migrations")
	}

	err = runMigrations(storage.DB)
	if err != nil {
		log.Fatal("failed to run migrations")
	}

	log.Println("migration process completed successfully")
}

func runMigrations(db *sql.DB) error {
	const op = "storage.postgresql.migrations.RunMigrations"

	data, err := os.ReadFile("data/data.json")
	if err != nil {
		return fmt.Errorf("%s: %w:", op, err)
	}

	var jsonData Data
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("%s: %w:", op, err)
	}

	for _, temp := range jsonData.Temperature {
		_, err := db.Exec("INSERT INTO temperature (temperature, recommendation)"+
			"VALUES ($1, $2)", temp.Temperature, temp.Recommendation)
		if err != nil {
			return fmt.Errorf("%s: %w:", op, err)
		}
	}

	for _, wind := range jsonData.Wind {
		_, err := db.Exec("INSERT INTO wind (wind_speed, recommendation)"+
			"VALUES ($1, $2", wind.WindSpeed, wind.Recommendation)
		if err != nil {
			return fmt.Errorf("%s: %w:", op, err)
		}
	}

	return nil
}
