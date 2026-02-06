package config

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Env string

const (
	EnvProduction Env = "production"
	EnvStaging    Env = "staging"
	EnvLocal      Env = "local"
)

type Config struct {
	TodoAddr     string       `required:"true" split_words:"true"`
	AllowOrigins AllowOrigins `split_words:"true"`
	// Environment
	Env Env `split_words:"true"`
	// Port for BFF
	Port string `split_words:"true"`
	// CORS headers to allow
	AllowHeaders []string `ignored:"true"`
}

type AllowOrigins []string

func Load() *Config {
	c := &Config{
		AllowHeaders: []string{
			"accept",
			"Accept-Encoding",
			"Accept-Language",
			"Connection",
			"Content-Length",
			"content-type",
			"Host",
			"Origin",
			"Referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"User-Agent",
			"x-access-token",
			"x-id-token",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Origin",
			"Authorization",
			"X-OS-TYPE",
			"X-OS-VERSION",
			"X-APP-VERSION",
		},
	}

	if err := envconfig.Process("", c); err != nil {
		panic(err)
	}

	return c
}

func (c *Config) AllowHosts() []string {
	res := make([]string, len(c.AllowOrigins))
	for i, s := range c.AllowOrigins {
		if s != "*" {
			parsedURL, err := url.Parse(s)
			if err != nil {
				panic(err)
			}
			res[i] = parsedURL.Host
		} else {
			res[i] = s
		}
	}
	return res
}
