package service

import (
	"go-api-samp/infrastructure"
	"go-api-samp/repository"
	service "go-api-samp/service/interface"
	"go-api-samp/service/internal/api"
)

/*
 * Service Provider
 */

type Provider interface {
	GetAPIService() service.APIService
}

type DefaultProvider struct {
	InfrastructureProvider infrastructure.Provider
	RepositoryProvider     repository.Provider
}

func (p *DefaultProvider) GetAPIService() service.APIService {
	return &api.API{
		Store: p.RepositoryProvider.GetWeatherStore(),
		Infra: p.InfrastructureProvider.GetMetaDataManager(),
	}
}
