package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Debug              bool
	GeniusRapidApiHost string        `split_words:"true" required:"true"`
	GeniusHost         string        `split_words:"true" required:"true"`
	GeniusRapidApiKey  string        `split_words:"true" required:"true"`
	GeniusApiHost      string        `split_words:"true" required:"true"`
	UserAgents         []string      `split_words:"true"`
	RequestTimeout     time.Duration `split_words:"true" default:"5s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	// TODO: Sprawdzanie, czy HOST kończy się na "/"
	return cfg, err
}
