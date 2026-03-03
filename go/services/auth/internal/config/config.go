package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env                  string `required:"true" default:"local"`
	ServerPort           int    `required:"true" split_words:"true"`
	GrpcReflectionEnable bool   `required:"true" split_words:"true"`
	TodoAddr             string `required:"true" split_words:"true"`

	DBConfig
	JwtConfig
}

func LoadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &c, nil
}
