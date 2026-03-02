package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type JwtConfig struct {
	AccessTokenSecret         string `required:"true" split_words:"true"`
	RefreshTokenSecret        string `required:"true" split_words:"true"`
	AccessTokenExpiresSecond  int    `required:"true" split_words:"true"`
	RefreshTokenExpiresSecond int    `required:"true" split_words:"true"`
}

func LoadJwtConfig() (*JwtConfig, error) {
	var c JwtConfig
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt config: %w", err)
	}

	return &c, nil
}
