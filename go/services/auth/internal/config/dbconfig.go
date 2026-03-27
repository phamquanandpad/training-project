package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	DBHost string `required:"true" split_words:"true"`
	DBPort int    `required:"true" split_words:"true"`
	DBUser string `required:"true" split_words:"true"`
	DBPass string `required:"true" split_words:"true"`
	DBName string `required:"true" split_words:"true"`
}

func LoadDBConfig() (*DBConfig, error) {
	var c DBConfig
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed to load db config: %w", err)
	}

	return &c, nil
}
