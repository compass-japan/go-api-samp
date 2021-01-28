package repository

import (
	"context"
	"go-api-samp/model/entity"
)

type (
	WeatherStoreManager interface {
		AddWeather(ctx context.Context, locationId, weather int, date, comment string) error
		GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error)
		GetLocation(ctx context.Context, locationId int) (bool, error)
	}
)
