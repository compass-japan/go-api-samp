package api

import (
	"context"
	infrastructure "go-api-samp/infrastructure/interface"
	"go-api-samp/model"
	"go-api-samp/model/dto"
	"go-api-samp/model/errors"
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

	_, err = a.Store.GetWeather(ctx, payload.LocationId, payload.Date)
	isUpdate := false
	if v, ok := err.(errors.SystemError); ok {
		switch {
		case v.Is(errors.DataStoreSystemError(nil)):
			logger.Error(ctx, "failed to get weather", err)
			return err
		case v.Is(errors.DataStoreValueNotFoundSystemError(nil)):
			isUpdate = true
		}
	}

	if isUpdate {
		if err := a.Store.UpdateWeather(ctx, payload.LocationId, payload.Weather, payload.Date, payload.Comment); err != nil {
			logger.Error(ctx, "failed to update weather.", err)
			return err
		}
	} else {
		if err := a.Store.AddWeather(ctx, payload.LocationId, payload.Weather, payload.Date, payload.Comment); err != nil {
			logger.Error(ctx, "failed to add weather.", err)
			return err
		}
	}

	return nil
}

func (a *API) GetWeather(ctx context.Context, payload *dto.GetWeatherRequest) (*dto.GetWeatherResponse, error) {
	logger := log.GetLogger()

	ety, err := a.Store.GetWeather(ctx, payload.LocationId, payload.Date)
	if err != nil {
		logger.Error(ctx, "failed to get weather", err)
		return nil, err
	}

	response := &dto.GetWeatherResponse{
		Location: ety.Location.City,
		Date:     ety.Dat,
		Weather:  model.ToWeather(ety.Weather),
		Comment:  ety.Comment,
	}

	return response, nil
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
