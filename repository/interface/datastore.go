package repository

import "context"

type (
	WeatherStoreManager interface {
		Add(ctx context.Context, locationId int, date, weather, comment string) error
		Get(ctx context.Context, locationId int, date string) (string, error)
	}
)
