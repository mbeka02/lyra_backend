package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SYMMETRIC_KEY                string        `mapstructure:"SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION        time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SENDGRID_API_KEY             string        `mapstructure:"SENDGRID_API_KEY"`
	GCLOUD_PROJECT_ID            string        `mapstructure:"GCLOUD_PROJECT_ID"`
	GCLOUD_IMAGE_BUCKET          string        `mapstructure:"GCLOUD_IMAGE_BUCKET"`
	GCLOUD_PATIENT_RECORD_BUCKET string        `mapstructure:"GCLOUD_PATIENT_RECORD_BUCKET"`
	GCLOUD_DATASET_LOCATION      string        `mapstructure:"GCLOUD_DATASET_LOCATION"`
	GCLOUD_DATASET_ID            string        `mapstructure:"GCLOUD_DATASET_ID"`
	GCLOUD_FHIR_STORE_ID         string        `mapstructure:"GCLOUD_FHIR_STORE_ID"`
	PAYSTACK_API_KEY             string        `mapstructure:"PAYSTACK_API_KEY"`
	GETSTREAM_API_KEY            string        `mapstructure:"GETSTREAM_API_KEY"`
	GETSTREAM_API_SECRET         string        `mapstructure:"GETSTREAM_API_SECRET"`
	// DB_CONNECTION_STRING  string        `mapstructure:"DB_CONNECTION_STRING"`
	PORT string `mapstructure:"PORT"`
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
