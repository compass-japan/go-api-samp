package config

type (
	WholeConfig struct {
		Logging LogConfig `yaml:"logging"`
	}

	LogConfig struct {
		Level string `yaml:"level"`
	}
)

var (
	Log = &LogConfig{}
)
