package main

import (
	"os"

	"fmt"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/internal"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
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
				Name:     "is-album",
				Usage:    "--is-album - will iterate thru album from which the song has came",
				Aliases:  []string{"ia"},
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			inputs["song"] = c.String("song")
			inputs["keyword"] = c.String("keyword")
			inputs["is-album"] = c.String("is-album")
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
	_, isAlbum := inputs["is-album"]

	if err != nil {
		log.Fatal(err)
	}

	mainLogger := log.New()
	mainLogger.SetLevel(log.DebugLevel)
	logger := log.NewEntry(mainLogger)

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Error("config creating error")
	}

	genius := internal.NewGeniusProvider(client, &cfg, logger.WithField("component", "genius_provider"))

	searchResults, err := genius.Search(songName)
	if err != nil {
		logger.WithError(err).Fatal("search error")
	}

	if isAlbum {
		artistID := searchResults[0].PrimaryArtist.ID
		_, _ = genius.FindSongsByArtistID(artistID)
	}

	lyrics, err := genius.GetLyricsFromPath(searchResults[0].LyricsPath)
	if err != nil {
		logger.WithError(err).Fatal("geting lyrics error")
	}

	count := strings.Count(strings.ToLower(lyrics), strings.ToLower(keyword))
	fmt.Printf("\"%s\" occurred %d times in \"%s\"\n", keyword, count, searchResults[0].FullTitle)
}
