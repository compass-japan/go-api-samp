package api

import (
	"context"
	infrastructure "go-api-samp/infrastructure/interface"
	"go-api-samp/model/dto"
	"go-api-samp/model/entity"
	repository "go-api-samp/repository/interface"
	"go-api-samp/util/log"
)

type API struct {
	Store repository.WeatherStoreManager
	Infra infrastructure.MetaWeatherManager
}

func (a *API) Register(ctx context.Context, payload *dto.RegisterRequest) error {
	logger := log.GetLogger()

	err := a.Store.FindLocation(ctx, payload.LocationId)
	if err != nil {
		logger.Error(ctx, "failed to find location", err)
		return err
	}

	if err := a.Store.AddWeather(ctx, payload.LocationId, payload.Weather, payload.Date, payload.Comment); err != nil {
		logger.Error(ctx, "failed to register.", err)
		return err
	}

	return nil
}

func (a *API) GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (*entity.Weather, error) {
	logger := log.GetLogger()

	entity, err := a.Store.GetWeather(ctx, payload.LocationId, payload.Date)
	if err != nil {
		logger.Error(ctx, "failed to get weather", err)
		return nil, err
	}

	return entity, nil
}

func (a *API) GetAPIData(ctx context.Context) (*dto.ExApiResponse, error) {
	logger := log.GetLogger()

	w, err := a.Infra.GetExWeather(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get ex weather.", err)
		return nil, err
	}

	return w, nil
}
