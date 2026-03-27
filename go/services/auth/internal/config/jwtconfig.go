package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type JwtConfig struct {
	AccessTokenSecret          string `required:"true" split_words:"true"`
	RefreshTokenSecret         string `required:"true" split_words:"true"`
	AccessTokenExpireDuration  int64  `required:"true" split_words:"true"`
	RefreshTokenExpireDuration int64  `required:"true" split_words:"true"`
}

func LoadJwtConfig() (*JwtConfig, error) {
	var c JwtConfig
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt config: %w", err)
	}

	return &c, nil
}
