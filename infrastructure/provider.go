package infrastructure

import (
	infrastructure "go-api-samp/infrastructure/interface"
	"go-api-samp/infrastructure/internal/exapi"
	"go-api-samp/util/config"
)

type Provider interface {
	GetMetaDataManager() infrastructure.MetaWeatherManager
}

type DefaultProvider struct{}

func (p *DefaultProvider) GetMetaDataManager() infrastructure.MetaWeatherManager {
	return &exapi.MetaWeatherClient{
		Config: config.ExAPI,
	}
}
