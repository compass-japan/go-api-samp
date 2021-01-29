package service

import (
	"context"
	"go-api-samp/model/dto"
	"go-api-samp/model/entity"
)

type APIService interface {
	Register(ctx context.Context, payload *dto.RegisterRequest) error
	GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (*entity.Weather, error)
	GetAPIData(ctx context.Context) (*dto.ExApiResponse, error)
}
