package main

import (
	"fmt"
	"github.com/marosiak/WordFinder/api"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
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

	app, err := api.NewAPI(
		fmt.Sprintf(":%d", cfg.ServerPort),
		api.NewGeniusAPI(&cfg, lyricsService, logger),
		api.NewDictionaryAPI(&cfg, logger),
	)
	if err != nil {
		logger.WithError(err).Fatal("cannot create API")
	}

	log.Fatal(app.Listen())
}
