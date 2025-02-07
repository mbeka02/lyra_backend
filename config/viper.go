package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SYMMETRIC_KEY         string        `mapstructure:"SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SENDGRID_API_KEY      string        `mapstructure:"SENDGRID_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	// tell Viper the location of the config file
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// read values
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return config, err
	}
	return config, nil
}
