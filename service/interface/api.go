package service

import (
	"context"
	"go-api-samp/model/dto"
)

/*
 * Service Interface
 */

type APIService interface {
	Register(ctx context.Context, payload *dto.RegisterRequest) error
	GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (*dto.GetWeatherResponse, error)
	GetAPIData(ctx context.Context) (*dto.ExApiResponse, error)
}
