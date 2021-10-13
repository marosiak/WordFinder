package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	searchEndpoint   = "search"
	songEndpoint     = "songs"
	artistEndpoint   = "artists"
	emptyLyricsErr   = errors.New("empty lyrics")
	blockedSelectors = []string{"script", "#onetrust-consent-sdk"}
)

const maxLyricsRetries = 2

type Song struct {
	ID        int
	Path      string
	FullTitle string `json:"full_title"`
}

type Artist struct {
	ID      int    `json:"id"`
	ApiPath string `json:"api_path"`
	Name    string `json:"name"`
}

type SearchResult struct {
	ID            int
	ApiPath       string `json:"api_path"`
	FullTitle     string `json:"full_title"`
	LyricsPath    string `json:"path"`
	PrimaryArtist Artist `json:"primary_artist"`
}

var _ GeniusProvider = &InternalGeniusProvider{}

type GeniusProvider interface {
	Search(query string) ([]SearchResult, error)
	GetSong(id int) (Song, error)
	GetLyrics(song Song) (string, error)
	GetLyricsFromPath(lyricsPath string) (string, error)
	FindSongsByArtistID(artistID int) (*Artist, error)
}

type InternalGeniusProvider struct {
	client *http.Client
	logger *log.Entry
	cfg    *config.Config
}

func NewGeniusProvider(client *http.Client, cfg *config.Config, logger *log.Entry) *InternalGeniusProvider {
	return &InternalGeniusProvider{client: client, cfg: cfg, logger: logger}
}

func (s *InternalGeniusProvider) GetSong(id int) (Song, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, fmt.Sprintf("%s/%d", songEndpoint, id), "GET")
	if err != nil {
		s.logger.WithError(err).Error("creating url")
		return Song{}, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		s.logger.WithError(err).Error("creating http client")
		return Song{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
	}

	type songResponse struct {
		Response struct {
			Song Song `json:"song"`
		}
	}

	songPayload := songResponse{}
	err = json.Unmarshal(bytes, &songPayload)
	return songPayload.Response.Song, err
}

func (s *InternalGeniusProvider) GetLyrics(song Song) (string, error) {
	return s.GetLyricsFromPath(song.Path)
}

func (s *InternalGeniusProvider) GetLyricsFromPath(lyricsPath string) (string, error) {
	req, err := utils.CreatePathRequest(s.cfg, fmt.Sprintf("%s%s", s.cfg.GeniusHost, lyricsPath), "GET")
	if err != nil {
		s.logger.WithError(err).Error("creating url")
		return "", err
	}

	retries := 1

REQUEST:
	res, err := s.client.Do(&req)
	if err != nil {
		s.logger.WithError(err).Error("creating http client")
		return "", err
	}

	if res.StatusCode != 200 {
		s.logger.Error(res.StatusCode)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.logger.WithError(err).WithFields(
			log.Fields{
				"status_code": res.StatusCode,
			}).Error("cannot query document")

		return "", err
	}

	// There is cookie popup and scripts which are hidding lyrics, this loop will remove it
	for _, v := range blockedSelectors {
		doc.Find(v).Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}

	var lyrics string
	doc.Find(".lyrics").Each(func(i int, s *goquery.Selection) {
		lyrics = s.Text()
	})

	if lyrics == "" {
		if retries < maxLyricsRetries {
			retries = retries + 1
			s.logger.WithError(emptyLyricsErr)
			goto REQUEST
		}

		return "", emptyLyricsErr
	}
	return lyrics, err
}
func (s *InternalGeniusProvider) Search(query string) ([]SearchResult, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, fmt.Sprintf("%s?q=%s", searchEndpoint, query), "GET")
	if err != nil {
		log.WithError(err).Error("creating url")
		return []SearchResult{}, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		log.WithError(err).Error("creating http client")
		return []SearchResult{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
	}

	type hit struct {
		Type         string
		SearchResult SearchResult `json:"result"`
	}

	type searchResultResponse struct {
		Response struct {
			Hits []hit `json:"hits"`
		}
	}

	searchResult := searchResultResponse{}
	err = json.Unmarshal(bytes, &searchResult)

	var results []SearchResult
	for _, hit := range searchResult.Response.Hits {
		results = append(results, hit.SearchResult)
	}

	return results, err
}

func (s *InternalGeniusProvider) FindSongsByArtistID(artistID int) (*Artist, error) {
	req, err := utils.CreatePathRequest(s.cfg, fmt.Sprintf("%s/%s/%d/%s", s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint), "GET")
	if err != nil {
		log.WithError(err).Error("creating url")
		return nil, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		log.WithError(err).Error("creating http client")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
	}
	println(string(bytes))
	return nil, errors.New("not implemented yet")
}
