package internal

import (
	"bufio"
	"github.com/marosiak/WordFinder/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
)

type Cmd interface {
	GetSongsByArtistWithoutBannedWords(ctx *cli.Context) error
}

var _ Cmd = &InternalCmd{}

type InternalCmd struct {
	lyricsService LyricsService
	logger        *log.Entry
	cfg           *config.Config
}

func NewCmd(cfg *config.Config, lyricsService LyricsService, logger *log.Entry) *InternalCmd {
	return &InternalCmd{logger: logger, lyricsService: lyricsService, cfg: cfg}
}

func getKeywordsFromFile(ctx *cli.Context) (output []string) {
	keywordsFile := ctx.String("keywords-file")
	f, err := os.OpenFile(keywordsFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		text := sc.Text() // GET the line string

		splittedSpace := strings.Split(text, " ")
		for _, word := range splittedSpace {
			splittedComma := strings.Split(word, ",")
			for _, word := range splittedComma {
				if word != "" {
					output = append(output, strings.ToLower(word))
				}
			}

		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
	return output
}

func getKeywords(ctx *cli.Context) (output []string) {
	keywords := ctx.StringSlice("keywords")
	for _, keyword := range keywords {
		splitted := strings.Split(keyword, ",")
		if len(splitted) > 1 {
			for _, v := range splitted {
				output = append(output, v)
			}
		} else {
			output = append(output, keyword)
		}
	}
	keywordsFromFile := getKeywordsFromFile(ctx)
	keywords = append(keywords, keywordsFromFile...)
	keywords = append(keywords, ctx.String("keyword"))
	return keywords
}

func (s *InternalCmd) GetSongsByArtistWithoutBannedWords(ctx *cli.Context) error {
	query := ctx.String("query")
	keywords := getKeywords(ctx)

	songs, err := s.lyricsService.GetSongsByArtist(query)
	if err != nil {
		return err
	}

	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Info.Title < songs[j].Info.Title
	})

	songToWordsMap := make(map[Song]WordsOccurrences)
	for _, song := range songs {
		songToWordsMap[song] = song.Lyrics.FindWords()
	}

	magicContainer := make(map[string]struct{})

	for _, keyword := range keywords {
		keyword = strings.ReplaceAll(keyword, " ", "")
		for song, occurrence := range songToWordsMap {
			if occurrence[keyword] >= 1 {

			} else {
				magicContainer[song.Info.Title] = struct{}{}
			}
		}
	}

	for k, _ := range magicContainer {
		println(k)
	}
	return nil
}
