package repository

import (
	"go-api-samp/application"
	repository "go-api-samp/repository/interface"
	"go-api-samp/repository/internal/datastore"
)

type (
	Provider interface {
		GetWeatherStore() repository.WeatherStoreManager
	}
	InitProvider interface {
		GetInitManager() repository.InitManager
	}
)

type DefaultProvider struct{}

func (p *DefaultProvider) GetWeatherStore() repository.WeatherStoreManager {
	return &datastore.MySQLClient{
		Db: application.GetDB(),
	}
}

type InitDefaultProvider struct{}

func (p *InitDefaultProvider) GetInitManager() repository.InitManager {
	return &datastore.InitClient{
		Db: application.GetDB(),
	}
}
