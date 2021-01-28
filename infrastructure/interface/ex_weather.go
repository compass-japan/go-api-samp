package infrastructure

import (
	"context"
	"go-api-samp/model/dto"
)

type (
	MetaWeatherManager interface {
		GetSample(ctx context.Context) (*dto.ExApiResponse, error)
	}
)
