package infrastructure

import (
	infrastructure "go-api-samp/infrastructure/interface"
	"go-api-samp/infrastructure/internal/ex"
)

type Provider interface {
	GetMetaDataManager() infrastructure.MetaWeatherManager
}

type DefaultProvider struct{}

func (p *DefaultProvider) GetMetaDataManager() infrastructure.MetaWeatherManager {
	return &ex.MetaWeatherClient{}
}
