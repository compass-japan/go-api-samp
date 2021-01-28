package repository

import (
	"go-api-samp/application"
	repository "go-api-samp/repository/interface"
	"go-api-samp/repository/internal/datastore"
)

type Provider interface {
	GetWeatherStore() repository.WeatherStoreManager
}

type DefaultProvider struct{}

func (p *DefaultProvider) GetWeatherStore() repository.WeatherStoreManager {
	return &datastore.MySQLClient{
		Db: application.GetDB(),
	}
}
