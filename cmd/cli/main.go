package main

import (
	"github.com/marosiak/WordFinder/utils"
	"os"

	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getCliApp() (map[string]string, error) {
	inputs := make(map[string]string)

	app := &cli.App{
		Name:  "genius-cli",
		Usage: "genius-cli --song=[song title] keyword=[keyword]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "query",
				Usage:    "--query=\"the_name\"",
				Aliases:  []string{"q"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "keyword",
				Usage:    "--keyword=\"the_keyword\"",
				Aliases:  []string{"k"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "scan-artist",
				Usage:    "--scan-artist - will iterate thru author of song provided with --song",
				Aliases:  []string{"ia"},
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			inputs["query"] = c.String("query")
			inputs["keyword"] = c.String("keyword")
			inputs["scan-artist"] = c.String("scan-artist")
			return nil
		},
	}

	err := app.Run(os.Args)
	return inputs, err
}

func main() {
	inputs, err := getCliApp()
	songName := inputs["query"]
	//keyword := inputs["keyword"]
	_, scanArtist := inputs["scan-artist"]

	if err != nil {
		log.Fatal(err)
	}

	mainLogger := log.New()
	mainLogger.SetLevel(log.DebugLevel)
	logger := log.NewEntry(mainLogger)

	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Error("config creating error")
	}

	genius := internal.NewGeniusProvider(utils.CreateHttpClient(&cfg), &cfg, logger.WithField("component", "genius_provider"))
	lyricsService := internal.NewLyricsService(&cfg, genius, logger)

	if scanArtist {
		songInfos, err := lyricsService.GetAllSongsInfoByArtist(songName)
		if err != nil {
			logger.WithError(err).Error("cannot fetch all songs infos by artist")
		}
		for _, v := range songInfos {
			println(v.Title)
		}

		song, err := lyricsService.GetSongFromInfo(songInfos[0])
		if err != nil {
			return
		}
		println(song.Lyrics)
	}

}
