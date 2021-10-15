package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Debug          bool
	GeniusApiHost  string        `split_words:"true"`
	GeniusHost     string        `split_words:"true"`
	GeniusApiKey   string        `split_words:"true"`
	UserAgents     []string      `split_words:"true"`
	RequestTimeout time.Duration `split_words:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	// TODO: Sprawdzanie, czy HOST kończy się na "/"
	return cfg, err
}
