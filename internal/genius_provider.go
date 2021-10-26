package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/marosiak/WordFinder/config"
	"github.com/marosiak/WordFinder/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

var (
	searchEndpoint   = "search"
	songEndpoint     = "songs"
	artistEndpoint   = "artists"
	emptyLyricsErr   = errors.New("empty lyrics")
	blockedSelectors = []string{"script", "#onetrust-consent-sdk"}
)

const (
	perPageLimit     = 50
	maxLyricsRetries = 2
)

type GeniusSong struct {
	ID            int
	LyricsPath    string       `json:"path"`
	FullTitle     string       `json:"full_title"`
	PrimaryArtist GeniusArtist `json:"primary_artist"`
	LyricsState   LyricsState  `json:"lyrics_state"`
}

type LyricsState string

const (
	LyricsComplete   LyricsState = "complete"
	LyricsIncomplete             = "incomplete"
	LyricsUnreleased             = "unreleased"
)

type GeniusArtist struct {
	ID      int    `json:"id"`
	ApiPath string `json:"api_path"`
	Name    string `json:"name"`
}

type SearchResult struct {
	ID             int
	ApiPath        string       `json:"api_path"`
	FullTitle      string       `json:"full_title"`
	LyricsEndpoint string       `json:"path"`
	PrimaryArtist  GeniusArtist `json:"primary_artist"`
}

var _ GeniusProvider = &InternalGeniusProvider{}

type GeniusProvider interface {
	Search(query string) ([]SearchResult, error)
	GetSongByID(id int) (GeniusSong, error)
	GetLyrics(song GeniusSong) (string, error)
	GetLyricsFromPath(lyricsPath string) (string, error)
	//FindSongsByArtistID(artistID int) ([]GeniusSong, error)
	FindSongsByArtistID(artistID int) ([]GeniusSong, error)
	FindArtist(artistName string) (GeniusArtist, error)
}

type InternalGeniusProvider struct {
	client *http.Client
	logger *log.Entry
	cfg    *config.Config
}

func NewGeniusProvider(client *http.Client, cfg *config.Config, logger *log.Entry) *InternalGeniusProvider {
	return &InternalGeniusProvider{client: client, cfg: cfg, logger: logger}
}

func (s *InternalGeniusProvider) FindArtist(artistName string) (GeniusArtist, error) {
	return GeniusArtist{}, nil
}
func (s *InternalGeniusProvider) GetSongByID(id int) (GeniusSong, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, fmt.Sprintf("%s/%d", songEndpoint, id), "GET")
	if err != nil {
		s.logger.WithError(err).Error("creating url")
		return GeniusSong{}, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		s.logger.WithError(err).Error("creating http client")
		return GeniusSong{}, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
	}

	type songResponse struct {
		Response struct {
			Song GeniusSong `json:"song"`
		}
	}

	songPayload := songResponse{}
	err = json.Unmarshal(bytes, &songPayload)
	return songPayload.Response.Song, err
}

func (s *InternalGeniusProvider) Search(query string) ([]SearchResult, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, fmt.Sprintf("%s?q=%s", searchEndpoint, url.QueryEscape(query)), "GET")
	if err != nil {
		log.WithError(err).Error("creating url")
		return []SearchResult{}, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		log.WithError(err).Error("creating http client")
		return []SearchResult{}, err
	}
	defer res.Body.Close()

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

//func (s *InternalGeniusProvider) FindSongsByArtistID(artistID int) ([]GeniusSong, error) {
//	var songs []GeniusSong
//	currentPage := 0
//
//	type artistSongsResponse struct {
//		Response struct {
//			Songs    []GeniusSong `json:"songs"`
//			NextPage int          `json:"next_page"`
//		} `json:"response"`
//	}
//
//REQUEST:
//	var url string
//	if currentPage == 0 {
//		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d", s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint, perPageLimit)
//	} else {
//		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d?page=%d", s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint, perPageLimit, currentPage)
//	}
//	println(url)
//	req, err := utils.CreatePathRequest(s.cfg, url, "GET")
//	if err != nil {
//		log.WithError(err).Error("creating url")
//		return nil, err
//	}
//
//	res, err := s.client.Do(&req)
//	if err != nil {
//		log.WithError(err).Error("creating http client")
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	by, err := io.ReadAll(res.Body)
//	if err != nil {
//		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
//		return nil, err
//	}
//
//	var artistSongsResp artistSongsResponse
//	err = json.Unmarshal(by, &artistSongsResp)
//
//	// The additional validation is needed, because sometimes the artist is on "feat" and the lyrics from feats aren't supported yet
//	for _, song := range artistSongsResp.Response.Songs {
//		if song.PrimaryArtist.ID == artistID {
//			songs = append(songs, song)
//		}
//	}
//
//	nextPage := artistSongsResp.Response.NextPage
//	if currentPage != nextPage {
//		currentPage = nextPage
//		goto REQUEST
//	}
//
//	return songs, err
//}

func (s *InternalGeniusProvider) FindSongsByArtistID(artistID int) ([]GeniusSong, error) {
	var songs []GeniusSong
	currentPage := 0

	type artistSongsResponse struct {
		Response struct {
			Songs    []GeniusSong `json:"songs"`
			NextPage int          `json:"next_page"`
		} `json:"response"`
	}

REQUEST:
	var url string
	if currentPage == 0 {
		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d", s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint, perPageLimit)
	} else {
		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d?page=%d", s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint, perPageLimit, currentPage)
	}
	println(url)
	req, err := utils.CreatePathRequest(s.cfg, url, "GET")
	if err != nil {
		log.WithError(err).Error("creating url")
		return nil, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		log.WithError(err).Error("creating http client")
		return nil, err
	}
	defer res.Body.Close()

	by, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
		return nil, err
	}

	var artistSongsResp artistSongsResponse
	err = json.Unmarshal(by, &artistSongsResp)

	// The additional validation is needed, because sometimes the artist is on "feat" and the lyrics from feats aren't supported yet
	for _, song := range artistSongsResp.Response.Songs {
		if song.PrimaryArtist.ID == artistID {
			songs = append(songs, song)
		}
	}

	nextPage := artistSongsResp.Response.NextPage
	if currentPage != nextPage {
		currentPage = nextPage
		goto REQUEST
	}

	return songs, err
}

func (s *InternalGeniusProvider) GetLyrics(song GeniusSong) (string, error) {
	return s.GetLyricsFromPath(song.LyricsPath)
}

func (s *InternalGeniusProvider) GetLyricsFromPath(lyricsPath string) (string, error) {
	urlStr := fmt.Sprintf("%s%s", s.cfg.GeniusHost, lyricsPath)
	req, err := utils.CreatePathRequest(s.cfg, urlStr, "GET")
	if err != nil {
		s.logger.WithError(err).Error("creating url")
		return "", err
	}

	s.logger.Info(urlStr)

	retries := 1
REQUEST:
	res, err := s.client.Do(&req)
	if err != nil {
		s.logger.WithError(err).Error("creating http client")
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		s.logger.WithField("status_code", res.StatusCode).Error("wrong status code")
	}

	buf, err := io.ReadAll(res.Body)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		s.logger.WithError(err).WithFields(
			log.Fields{
				"status_code": res.StatusCode,
			}).Error("cannot query document")

		return "", err
	}

	// There is cookie popup and scripts which are hiding lyrics, this loop will remove it
	for _, v := range blockedSelectors {
		doc.Find(v).Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}

	var lyrics string

	// This list of selectors is needed in order to work around AB Tests, in future there could be pattern scanning
	whitelistedSelectors := []string{"#lyrics-root-pin-spacer", ".lyrics"}
	for _, selector := range whitelistedSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			lyrics = s.Text()
		})
	}

	if lyrics == "" {
		s.logger.WithError(emptyLyricsErr)

		if retries < maxLyricsRetries {
			retries = retries + 1
			goto REQUEST
		}

		return "", emptyLyricsErr
	}
	return lyrics, err
}
