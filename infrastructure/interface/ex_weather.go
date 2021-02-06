package infrastructure

import (
	"context"
	"go-api-samp/model/dto"
)

/*
 * infrastructure(外部APIアクセス)のinterface
 */

type (
	MetaWeatherManager interface {
		GetExWeather(ctx context.Context) (*dto.ExApiResponse, error)
	}
)
