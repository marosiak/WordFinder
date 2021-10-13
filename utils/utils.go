package utils

import (
	"fmt"
	"github.com/marosiak/WordFinder/config"
	"net/http"
	"net/url"
)

func getHeaders(cfg *config.Config) map[string][]string {
	return map[string][]string{
		"Content-Type":    {"application/json"},
		"x-rapidapi-host": {cfg.GeniusApiHost},
		"x-rapidapi-key":  {cfg.GeniusApiKey},
		"User-Agent":      {"Mozilla/5.0 (X11; U; Linux is686; pl-PL; rv:1.7.10) Gecko/20050717 Firefox/1.0.6"},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
		"Set-Cookie":      {"CONSTANT=YES+shp.gws-20210701-0-RC1.pl+FX+631;AMP-CONSENT=amp-8QWzAroGD8LPo0rQpgV1-w"},
	}
}

func CreateEndpointRequest(cfg *config.Config, endpoint string, method string) (http.Request, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("https://%s/%s", cfg.GeniusApiHost, endpoint))
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

func CreatePathRequest(cfg *config.Config, path string, method string) (http.Request, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("https://%s", path))
	println(reqUrl.String())
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
