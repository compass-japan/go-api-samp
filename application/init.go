package application

import (
	"go-api-samp/util/config"
	"go-api-samp/util/log"
)

func Init() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	if err := NewDBOpen(config.DB); err != nil {
		//return err
	}

	log.NewLogger(config.Log)

	return nil
}
