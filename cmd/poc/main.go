package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

var blockedSelectors = []string{"script", "#onetrust-consent-sdk"}

func getHeaders() map[string][]string {
	return map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"*/*"},
		"Connection":    {"keep-alive"},
		"Cache-Control": {"no-cache"},
	}
}

func getLyrics(client *http.Client) string {
	u, _ := url.Parse("https://genius.com/Mata-100-dni-do-matury-lyrics")
	req := http.Request{
		Method: "GET",
		URL:    u,
		Header: getHeaders(),
	}

	res, err := client.Do(&req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range blockedSelectors {
		doc.Find(v).Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}

	var lyrics string
	whitelistedSelectors := []string{"#lyrics-root-pin-spacer", ".lyrics"}

	for _, selector := range whitelistedSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			lyrics = lyrics + s.Text()
		})
		if lyrics != "" {
			break
		}
	}

	return lyrics
}

func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     time.Second * 5,
		},
	}

	for {
		lyrics := getLyrics(client)
		if lyrics == "" {
			log.Fatal(":/")
		} else {
			log.Info("Success")
		}
	}
}
