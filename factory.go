package main

type ProviderFactory interface {
}

var factoryInstance ProviderFactory

func GetProviderFactory() ProviderFactory {
	if factoryInstance == nil {
		factoryInstance = &defaultProviderFactory{}
	}
	return factoryInstance
}

type defaultProviderFactory struct{}
