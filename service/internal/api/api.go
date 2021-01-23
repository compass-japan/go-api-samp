package api

import (
	"context"
	infrastructure "go-api-samp/infrastructure/interface"
	"go-api-samp/model/dto"
	repository "go-api-samp/repository/interface"
)

type API struct {
	Store repository.WeatherStoreManager
	Infra infrastructure.MetaWeatherManager
}

func (a *API) Register(ctx context.Context, payload *dto.RegisterRequest) error {
	return nil
}

func (a *API) GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (string, error) {
	return "", nil
}

func (a *API) GetAPIData(ctx context.Context) (string, error) {
	return "", nil
}
