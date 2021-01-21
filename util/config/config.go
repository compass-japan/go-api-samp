package config

import (
	"github.com/jinzhu/configor"
	"os"
)

type (
	WholeConfig struct {
		Logging LogConfig    `yaml:"logging"`
		Server  ServerConfig `yaml:"server"`
	}

	ServerConfig struct {
		Addr string `yaml:"addr"`
	}
	LogConfig struct {
		Level string `yaml:"level"`
	}
)

var (
	Server = &ServerConfig{}
	Log    = &LogConfig{}
)

const confPath = "config/config.yaml"

func LoadConfig() error {
	env := os.Getenv("env")

	whole := &WholeConfig{}
	if err := configor.New(&configor.Config{Environment: env, ENVPrefix: "API", Verbose: true}).Load(whole, confPath); err != nil {
		return err
	}

	Log = &whole.Logging
	Server = &whole.Server
	return nil
}
