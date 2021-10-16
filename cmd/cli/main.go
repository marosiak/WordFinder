package main

import (
	"github.com/marosiak/WordFinder/utils"
	"os"

	"fmt"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

func getCliApp() (map[string]string, error) {
	inputs := make(map[string]string)

	app := &cli.App{
		Name:  "genius-cli",
		Usage: "genius-cli --song=[song title] keyword=[keyword]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "song",
				Usage:    "--song=\"the_name\"",
				Aliases:  []string{"s"},
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
			inputs["song"] = c.String("song")
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
	songName := inputs["song"]
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
		logger.WithError(err).Error("config creating error")
	}

	genius := internal.NewGeniusProvider(utils.CreateHttpClient(&cfg), &cfg, logger.WithField("component", "genius_provider"))

	searchResults, err := genius.Search(songName)
	if err != nil {
		logger.WithError(err).Fatal("search error")
	}

	if scanArtist {
		artistID := searchResults[0].PrimaryArtist.ID

		songs, err := genius.FindSongsByArtistID(artistID)
		if err != nil {
			logger.WithError(err).Fatal("geting songs")
		}

		for _, song := range songs {
			song := song
			c := make(chan string)
			go func(c chan string) {
				lyrics, err := genius.GetLyricsFromPath(song.LyricsPath)
				if err != nil {
					logger.WithError(err).Fatal("geting lyrics error")
				}
				count := strings.Count(strings.ToLower(lyrics), strings.ToLower(keyword))
				c <- fmt.Sprintf("\"%s\" occurred %d times in \"%s\"\n", keyword, count, song.FullTitle)
			}(c)

			a := <-c
			println(a)
		}

		return
	}

	lyrics, err := genius.GetLyricsFromPath(searchResults[0].LyricsPath)
	if err != nil {
		logger.WithError(err).Fatal("geting lyrics error")
	}

	count := strings.Count(strings.ToLower(lyrics), strings.ToLower(keyword))
	fmt.Printf("\"%s\" occurred %d times in \"%s\"\n", keyword, count, searchResults[0].FullTitle)
}
