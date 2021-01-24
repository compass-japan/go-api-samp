package application

import "go-api-samp/util/config"

func Init() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	if err := NewDBOpen(config.DB); err != nil {
		return err
	}

	return nil
}
