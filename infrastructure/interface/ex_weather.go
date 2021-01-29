package infrastructure

import (
	"context"
	"go-api-samp/model/dto"
)

type (
	MetaWeatherManager interface {
		GetExWeather(ctx context.Context) (*dto.ExApiResponse, error)
	}
)
