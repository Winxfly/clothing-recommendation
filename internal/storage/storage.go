package storage

import "errors"

var (
	ErrTemperatureNotFound = errors.New("recommendation for temperature not found")
	ErrWindNotFound        = errors.New("recommendation for wind not found")
)
