package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Debug                bool
	GeniusRapidApiHost   string        `split_words:"true" required:"true"`
	GeniusHost           string        `split_words:"true" required:"true"`
	GeniusRapidApiKey    string        `split_words:"true" required:"true"`
	GeniusApiHost        string        `split_words:"true" required:"true"`
	UserAgents           []string      `split_words:"true"`
	RequestTimeout       time.Duration `split_words:"true" default:"5s"`
	MaxChannelBufferSize int           `split_words:"true" default:"30"`
	MaxPagesForArtist    int           `split_words:"true" default:"100"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)

	return cfg, err
}
