package main

import (
	"fmt"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"time"
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
	startTime := time.Now()
	inputs, err := getCliApp()
	query := inputs["query"]
	keyword := inputs["keyword"]
	_, scanArtist := inputs["scan-artist"]

	if err != nil {
		log.Fatal(err)
	}

	mainLogger := log.New()
	mainLogger.SetLevel(log.DebugLevel)
	logger := log.NewEntry(mainLogger)

	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Fatal("config creating error")
	}

	genius := internal.NewGeniusProvider(utils.CreateHttpClient(&cfg), &cfg, logger.WithField("component", "genius_provider"))
	lyricsService := internal.NewLyricsService(&cfg, genius, logger)
	if scanArtist {
		songs, err := lyricsService.GetAllSongsByArtist(query)
		if err != nil {
			logger.WithError(err).Fatal("cannot fetch all songs by artist")

		}

		occurredAtleastOnceCounter := 0
		results := make(map[string]int)
		for _, song := range songs {
			results[song.Info.Title] = strings.Count(strings.ToLower(song.Lyrics), strings.ToLower(keyword))
		}

		for _, val := range results {
			if val > 0 {
				occurredAtleastOnceCounter = occurredAtleastOnceCounter + 1
			}
		}

		println("\n\n\n")
		for key, val := range results {
			fmt.Printf("%d times : \"%s\"\n", val, key)
		}

		println("\n================================\n")
		fmt.Printf("\nWord \"%s\" occurred in %d out of %d songs by %s\n", keyword, occurredAtleastOnceCounter, len(results), songs[0].Info.AuthorName)
	}
	fmt.Printf("\nCli ended after %s", time.Since(startTime).String())
}
