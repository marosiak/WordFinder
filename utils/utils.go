package utils

import (
	"fmt"
	"github.com/marosiak/WordFinder/config"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func getHeaders(cfg *config.Config) map[string][]string {
	rand.Seed(time.Now().UnixNano())
	randomUserAgent := cfg.UserAgents[rand.Intn(len(cfg.UserAgents)-0)]

	return map[string][]string{
		"Content-Type":    {"application/json"},
		"x-rapidapi-host": {cfg.GeniusRapidApiHost},
		"x-rapidapi-key":  {cfg.GeniusRapidApiKey},
		"User-Agent":      {randomUserAgent},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
		"Cache-Control":   {"no-cache"},
	}
}

func CreateEndpointRequest(cfg *config.Config, host string, endpoint string, method string) (http.Request, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("https://%s/%s", host, endpoint))
	if err != nil {
		return http.Request{}, err
	}

	if cfg.Debug {
		println(reqUrl.String())
	}

	req := http.Request{
		Method: method,
		URL:    reqUrl,
		Header: getHeaders(cfg),
	}
	return req, nil
}

func CreatePathRequest(cfg *config.Config, path string, method string) (http.Request, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("https://%s", path))
	if cfg.Debug {
		println(reqUrl.String())
	}

	if err != nil {
		return http.Request{}, err
	}

	req := http.Request{
		Method: method,
		URL:    reqUrl,
		Header: getHeaders(cfg),
	}
	return req, nil
}

func CreateHttpClient(cfg *config.Config) *http.Client {
	return &http.Client{
		Timeout: cfg.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     cfg.RequestTimeout,
		},
	}
}
