package internal

import (
	"encoding/base64"
	"fmt"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/response"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strings"
)

type GeniusAPI interface {
	GetSongsByArtist(ctx *fasthttp.RequestCtx)
	GetSongsWithWordsByArtist(ctx *fasthttp.RequestCtx)
}

type InternalGeniusAPI struct {
	lyricsService LyricsService
	cfg           *config.Config
	logger        *log.Entry
}

func NewGeniusAPI(cfg *config.Config, lyricsService LyricsService, logger *log.Entry) *InternalGeniusAPI {
	return &InternalGeniusAPI{cfg: cfg, lyricsService: lyricsService, logger: logger}
}

type apiSong struct {
	Title      string           `json:"title"`
	URL        string           `json:"url"`
	WordsCount WordsOccurrences `json:"words_count,omitempty"`
}

func (s *InternalGeniusAPI) GetSongsByArtist(ctx *fasthttp.RequestCtx) {
	type responseStruct struct {
		Songs []apiSong `json:"songs"`
	}

	artistName := ctx.Value("artist_name").(string)

	songs, err := s.lyricsService.GetSongsInfosByArtist(artistName)
	if err != nil {
		s.logger.WithError(err).Error("error getting songs infos by artist")
		response.WriteError(ctx, response.ErrorByName("internal_error"))
		return
	}

	resp := responseStruct{}
	for _, song := range songs {
		resp.Songs = append(resp.Songs, apiSong{
			Title: song.Title,
			URL:   fmt.Sprintf("https://%s%s", s.cfg.GeniusHost, song.PageEndpoint),
		})
	}
	response.WriteJSON(ctx, 200, response.New{
		Data: resp,
	})
}

type BannedWords []Word

func (b BannedWords) Normalise() BannedWords {
	var outputWords BannedWords
	for _, bannedWord := range b {
		for _, bannedChar := range "ĄĆĘŁŃÓŚŹŻ" {
			outputWords = append(outputWords, Word(strings.ReplaceAll(string(bannedWord), string(bannedChar), "")))
			outputWords = append(outputWords, Word(strings.ReplaceAll(string(bannedWord), strings.ToLower(string(bannedChar)), "")))
		}
	}
	return outputWords
}
func (b BannedWords) Contains(text Word) bool {
	for _, v := range b {
		if v == text {
			return true
		}
	}
	return false
}
func QueryStringList(ctx *fasthttp.RequestCtx, name string) BannedWords {
	by, _ := base64.StdEncoding.DecodeString(string(ctx.QueryArgs().Peek(name)))
	var bannedWords BannedWords
	for _, word := range strings.Split(string(by), ",") {
		bannedWords = append(bannedWords, Word(word))
	}
	return bannedWords
}

func (s *InternalGeniusAPI) GetSongsWithWordsByArtist(ctx *fasthttp.RequestCtx) {
	type responseStruct struct {
		Songs []apiSong
	}

	artistName := ctx.Value("artist_name").(string)
	bannedWords := QueryStringList(ctx, "banned_words")

	songs, err := s.lyricsService.GetSongsByArtist(artistName)
	if err != nil {
		s.logger.WithError(err).Error("error getting songs infos by artist")
		response.WriteError(ctx, response.ErrorByName("internal_error"))
		return
	}

	resp := responseStruct{}
	for _, song := range songs {
		if song.Lyrics.FindWords().ContainsOneOfWords(bannedWords.Normalise()) == false {
			resp.Songs = append(resp.Songs, apiSong{
				Title:      song.Info.Title,
				URL:        fmt.Sprintf("https://%s%s", s.cfg.GeniusHost, song.Info.PageEndpoint),
				WordsCount: song.Lyrics.FindWords(),
			})
		}
	}
	response.WriteJSON(ctx, 200, response.New{Data: resp})
}
