package repository

import (
	"context"
	"go-api-samp/model/entity"
)

type (
	WeatherStoreManager interface {
		AddWeather(ctx context.Context, locationId, weather int, date, comment string) error
		UpdateWeather(ctx context.Context, locationId, weather int, date, comment string) error
		GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error)
		FindLocation(ctx context.Context, locationId int) error
	}

	InitManager interface {
		FindAllLocations() ([]entity.Location, error)
	}
)
