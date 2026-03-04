package testutil

import (
	"github.com/kelseyhightower/envconfig"
)

type TestEnv struct {
	DBHost string `envconfig:"DB_HOST"`
	DBPort int    `envconfig:"DB_PORT"`
	DBUser string `envconfig:"DB_USER"`
	DBPass string `envconfig:"DB_PASS"`
	DBName string `envconfig:"DB_NAME"`
}

func LoadEnv() *TestEnv {
	var env TestEnv
	if err := envconfig.Process("", &env); err != nil {
		panic(err)
	}
	return &env
}
