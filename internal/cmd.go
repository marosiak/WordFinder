package internal

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/marosiak/WordFinder/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
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

func getKeywordsFromFile(fileName string) (output []string) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
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

var CannotOpenFileError = errors.New("Cannot open \"%s\" file\n")

func getKeywords(ctx *cli.Context) (output []string) {
	keywords := ctx.StringSlice("keywords")
	for _, keyword := range keywords {
		splitted := strings.Split(keyword, ",")
		for _, v := range splitted {
			output = append(output, v)
		}
	}

	keywordsFiles := ctx.StringSlice("keywords-files")
	keywordsFromFiles := []string{}
	for _, keyword := range keywordsFiles {
		splitted := strings.Split(keyword, ",")
		for _, fileName := range splitted {
			keywordsFromFiles = append(keywordsFromFiles, fileName)

			kw := getKeywordsFromFile(fileName)
			if len(kw) == 0 {
				fmt.Printf(CannotOpenFileError.Error(), fileName)
				continue
			}
			output = append(output, kw...)
		}
	}

	filePath := ctx.String("keywords-file")
	if filePath != "" {
		kw := getKeywordsFromFile(filePath)
		if len(kw) == 0 {
			fmt.Printf(CannotOpenFileError.Error(), filePath)
		}
		output = append(output, kw...)
	}

	keyword := ctx.String("keyword")
	if keyword != "" {
		output = append(output, keyword)
	}

	return output
}

func (s *InternalCmd) GetSongsByArtistWithoutBannedWords(ctx *cli.Context) error {
	query := ctx.String("query")
	keywords := getKeywords(ctx)

	songs, err := s.lyricsService.GetSongsByArtist(query)
	if err != nil {
		fmt.Printf("Error while getting songs list: %v", err)
		return err
	}

	songToWordsMap := make(map[Song]WordsOccurrences)
	for _, song := range songs {
		songToWordsMap[song] = song.Lyrics.FindWords()
	}

	songsWithoutBannedWords := make(map[string]struct{})
	for song, occurrence := range songToWordsMap {
		keywordExists := false
		for _, keyword := range keywords {
			keyword = strings.ReplaceAll(keyword, " ", "")
			if occurrence[keyword] > 0 {
				keywordExists = true
			}
		}
		if keywordExists == false {
			songsWithoutBannedWords[song.Info.Title] = struct{}{}
		}
	}

	for k, _ := range songsWithoutBannedWords {
		fmt.Println(k)
	}
	return nil
}
