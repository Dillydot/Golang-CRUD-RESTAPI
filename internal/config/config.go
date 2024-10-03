package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Postgres Postgres
}

type Postgres struct {
	Username string `mapstructure:"POSTGRES_USERNAME"`
	Pass     string `mapstructure:"POSTGRES_PASSWORD"`
	Host     string `mapstructure:"POSTGRES_HOSTNAME"`
	Database string `mapstructure:"POSTGRES_DATABASE"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
}

// LoadConfig loads configuration values from a file or env vars.
func LoadConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	viper.AutomaticEnv()

	var p Postgres
	err = viper.Unmarshal(&p)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Postgres: p,
	}, nil
}
