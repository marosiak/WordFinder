package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/marosiak/WordFinder/api"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	mainLogger := log.New()
	logger := log.NewEntry(mainLogger)
	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Fatal("config creating error")
	}

	if cfg.Debug {
		mainLogger.SetLevel(log.DebugLevel)
	} else {
		mainLogger.SetLevel(log.WarnLevel)
	}

	geniusProvider := internal.NewGeniusProvider(utils.CreateHttpClient(&cfg), &cfg, logger)
	lyricsService := internal.NewLyricsService(&cfg, geniusProvider, logger)

	geniusAPI := api.NewGeniusAPI(&cfg, lyricsService, logger)

	r := fasthttprouter.New()

	r.GET("/artists/:artist_name/songs", geniusAPI.GetSongsByArtist)
	r.GET("/artists/:artist_name/songs/words", geniusAPI.GetSongsWithWordsByArtist)

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r.Handler))
}
