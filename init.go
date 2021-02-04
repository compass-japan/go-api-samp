package main

import (
	"go-api-samp/application"
	"go-api-samp/repository"
	repositoryif "go-api-samp/repository/interface"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
)

func Init(provider repository.InitProvider) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	if err := application.NewDBOpen(config.DB); err != nil {
		//return err
	}

	if err := loadLocations(provider.GetInitManager()); err != nil {
		//return err
	}

	log.NewLogger(config.Log)

	return nil
}

func loadLocations(manager repositoryif.InitManager) error {
	locations, err := manager.FindAllLocations()
	if err != nil {
		return err
	}

	for _, v := range *locations {
		application.GetLocationsMap()[v.Id] = v.City
	}

	return nil
}
