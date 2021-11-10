package main

import (
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
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

	cmd := internal.NewCmd(&cfg, lyricsService, logger)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "query",
			Usage:    "--query=\"the_name\"",
			Aliases:  []string{"q"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keyword",
			Usage:    "--keyword=\"the_keyword\"",
			Aliases:  []string{"kwd"},
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "keywords",
			Usage:    "--keywords=\"the_keyword\",\"another_keyword\"",
			Aliases:  []string{"kwds"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "keywords-file",
			Usage:    "--keywords-file=\"keywords.txt\" the keywords should be splitted by space or \",\"",
			Aliases:  []string{"kwds-f"},
			Required: false,
		},
	}

	app := &cli.App{
		Name:  "genius-cli",
		Usage: "genius-cli --help",
		Commands: []*cli.Command{
			{
				Name:   "songs-by-artist-without-banned-words", // damnn.. I have to work it around
				Usage:  "Will return list of songs which does not contains any of --keywords or --keyword",
				Action: cmd.GetSongsByArtistWithoutBannedWords,
				Flags:  flags,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
