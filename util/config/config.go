package config

import (
	"github.com/jinzhu/configor"
	"os"
)

type (
	WholeConfig struct {
		Server  ServerConfig `yaml:"server"`
		Logging LogConfig    `yaml:"logging"`
		DB      DBConfig     `yaml:"db"`
	}

	ServerConfig struct {
		Addr string `yaml:"addr"`
	}
	LogConfig struct {
		Level string `yaml:"level"`
	}
	DBConfig struct {
		DriverName string `yaml:"driverName"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		DBName     string `yaml:"dbname"`
		Addr       string `yaml:"addr"`
	}
)

var (
	Server = &ServerConfig{}
	Log    = &LogConfig{}
	DB     = &DBConfig{}
)

const confPath = "config/config.yaml"

func LoadConfig() error {
	env := os.Getenv("env")

	whole := &WholeConfig{}
	if err := configor.New(&configor.Config{Environment: env, ENVPrefix: "API", Verbose: true}).Load(whole, confPath); err != nil {
		return err
	}

	Server = &whole.Server
	Log = &whole.Logging
	DB = &whole.DB
	return nil
}
