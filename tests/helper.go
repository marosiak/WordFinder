package tests

import "github.com/marosiak/WordFinder/config"

func GetConfig() *config.Config {
	return &config.Config{
		Debug:              true,
		GeniusRapidApiHost: "example.com",
		GeniusHost:         "example.com",
		GeniusRapidApiKey:  "api_key",
		GeniusApiHost:      "example.com",
		RequestTimeout:     5,
	}
}
