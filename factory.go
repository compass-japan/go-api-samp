package main

import (
	"go-api-samp/infrastructure"
	"go-api-samp/repository"
	"go-api-samp/service"
)

type ProviderFactory interface {
	GetServiceProvider() service.Provider
	GetInitProvider() repository.InitProvider
}

var factoryInstance ProviderFactory

func GetProviderFactory() ProviderFactory {
	if factoryInstance == nil {
		factoryInstance = &defaultProviderFactory{}
	}
	return factoryInstance
}

type defaultProviderFactory struct{}

func (p *defaultProviderFactory) GetServiceProvider() service.Provider {
	return &service.DefaultProvider{
		InfrastructureProvider: &infrastructure.DefaultProvider{},
		RepositoryProvider:     &repository.DefaultProvider{},
	}
}

func (p *defaultProviderFactory) GetInitProvider() repository.InitProvider {
	return &repository.InitDefaultProvider{}
}
