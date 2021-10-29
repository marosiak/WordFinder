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
	"sync"
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

type GeniusSongInfo struct {
	ID            int
	LyricsPath    string       `json:"path"`
	FullTitle     string       `json:"full_title"`
	PrimaryArtist GeniusArtist `json:"primary_artist"`
	LyricsState   LyricsState  `json:"lyrics_state"`
}

type GeniusSong struct {
	Lyrics string
	Info   GeniusSongInfo
}

type geniusSonginfos []GeniusSongInfo

func (s geniusSonginfos) ExistsByID(id int) bool {
	for _, song := range s {
		if song.ID == id {
			return true
		}
	}
	return false
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
	LyricsState    string       `json:"lyrics_state"`
}

var _ GeniusProvider = &InternalGeniusProvider{}

type GeniusProvider interface {
	Search(query string) ([]SearchResult, error)
	GetSongInfoByID(id int) (GeniusSongInfo, error)
	GetSongByID(id int) (GeniusSong, error)
	GetSongsByIDs(id []int) ([]GeniusSong, error)
	GetSongByName(name string) (GeniusSong, error)
	FindArtist(artistName string) (GeniusArtist, error)
	FindSongInfosByArtistID(artistID int) ([]GeniusSongInfo, error)
	FindSongsByArtistID(artistID int) ([]GeniusSong, error)
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
	songInfo, err := s.GetSongInfoByID(id)
	if err != nil {
		s.logger.WithError(err).Error("GetSongByID getting song info by ID ", id)
	}

	lyrics, err := s.getLyrics(songInfo)
	return GeniusSong{
		Lyrics: lyrics,
		Info:   songInfo,
	}, err
}

func (s *InternalGeniusProvider) GetSongsByIDs(ids []int) ([]GeniusSong, error) {
	songCh := make(chan GeniusSong, s.cfg.MaxChannelBufferSize)
	wg := sync.WaitGroup{}

	for _, id := range ids {
		wg.Add(1)
		id := id
		go func() {
			defer wg.Done()
			song, err := s.GetSongByID(id)
			if err != nil {
				s.logger.WithError(err).Error()
				return
			}
			songCh <- song
		}()

	}

	go func() {
		wg.Wait()
		close(songCh)
	}()

	var songs []GeniusSong
	for song := range songCh {
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *InternalGeniusProvider) GetSongByName(name string) (GeniusSong, error) {
	searchResults, err := s.Search(name)
	if err != nil {
		return GeniusSong{}, err
	}

	songResult := searchResults[0]
	artist := songResult.PrimaryArtist
	lyrics, err := s.getLyricsFromPath(songResult.LyricsEndpoint)

	return GeniusSong{
		Lyrics: lyrics,
		Info: GeniusSongInfo{
			ID:         songResult.ID,
			LyricsPath: songResult.LyricsEndpoint,
			FullTitle:  songResult.FullTitle,
			PrimaryArtist: GeniusArtist{
				ID:      artist.ID,
				ApiPath: artist.ApiPath,
				Name:    artist.Name,
			},
			LyricsState: LyricsState(songResult.LyricsState),
		},
	}, err
}

func (s *InternalGeniusProvider) GetSongInfoByID(id int) (GeniusSongInfo, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, s.cfg.GeniusRapidApiHost, fmt.Sprintf("%s/%d", songEndpoint, id), "GET")
	if err != nil {
		s.logger.WithError(err).Error("creating url")
		return GeniusSongInfo{}, err
	}

	res, err := s.client.Do(&req)
	if err != nil {
		s.logger.WithError(err).Error("creating http client")
		return GeniusSongInfo{}, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
	}

	type songResponse struct {
		Response struct {
			Song GeniusSongInfo `json:"song"`
		}
	}

	songPayload := songResponse{}
	err = json.Unmarshal(bytes, &songPayload)
	return songPayload.Response.Song, err
}

func (s *InternalGeniusProvider) Search(query string) ([]SearchResult, error) {
	req, err := utils.CreateEndpointRequest(s.cfg, s.cfg.GeniusRapidApiHost, fmt.Sprintf("%s?q=%s", searchEndpoint, url.QueryEscape(query)), "GET")
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

func getUrlForPage(host string, artistEndpoint string, artistID int, songEndpoint string, page int) string {
	var url string

	if page == 1 {
		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d", host, artistEndpoint, artistID, songEndpoint, perPageLimit)
	} else if page >= 3 {
		url = fmt.Sprintf("%s/%s/%d/%s?page=%d", host, artistEndpoint, artistID, songEndpoint, page)
	} else {
		url = fmt.Sprintf("%s/%s/%d/%s?per_page=%d?page=%d", host, artistEndpoint, artistID, songEndpoint, perPageLimit, page)
	}
	return url
}

func (s *InternalGeniusProvider) FindSongsByArtistID(artistID int) ([]GeniusSong, error) {
	songInfos, err := s.FindSongInfosByArtistID(artistID)
	if err != nil {
		return []GeniusSong{}, err
	}

	songCh := make(chan GeniusSong, s.cfg.MaxChannelBufferSize)
	wg := sync.WaitGroup{}

	for _, songInfo := range songInfos {
		wg.Add(1)
		songInfo := songInfo
		go func() {
			defer wg.Done()

			lyrics, err := s.getLyrics(songInfo)
			if err != nil {
				return
			}

			songCh <- GeniusSong{
				Lyrics: lyrics,
				Info:   songInfo,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(songCh)
	}()

	var songs []GeniusSong
	for song := range songCh {
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *InternalGeniusProvider) FindSongInfosByArtistID(artistID int) ([]GeniusSongInfo, error) {
	var songs geniusSonginfos

	type artistSongsResponse struct {
		Response struct {
			Songs    []GeniusSongInfo `json:"songs"`
			NextPage int              `json:"next_page"`
		} `json:"response"`
	}

	artistSongRespCh := make(chan artistSongsResponse, s.cfg.MaxChannelBufferSize)
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		i := i
		go func() {
			url := getUrlForPage(s.cfg.GeniusApiHost, artistEndpoint, artistID, songEndpoint, i)
			if s.cfg.Debug {
				println(url)
			}

			wg.Add(1)
			defer wg.Done()

			req, err := utils.CreatePathRequest(s.cfg, url, "GET")
			if err != nil {
				log.WithError(err).Error("creating url")
				errCh <- err
				return
			}

			res, err := s.client.Do(&req)
			if err != nil {
				log.WithError(err).Error("creating http client")
				errCh <- err
				return
			}
			defer res.Body.Close()

			by, err := io.ReadAll(res.Body)
			if err != nil {
				s.logger.WithError(err).Errorf("reading response status: %s", res.Status)
				errCh <- err
				return
			}

			var artistSongsResp artistSongsResponse
			err = json.Unmarshal(by, &artistSongsResp)
			artistSongRespCh <- artistSongsResp
		}()
	}

	go func() {
		wg.Wait()
		close(artistSongRespCh)
		close(errCh)
	}()

	// The additional validation is needed, because sometimes the artist is on "feat" and the lyrics from feats aren't supported yet
	for song := range artistSongRespCh {
		for _, song := range song.Response.Songs {
			if song.PrimaryArtist.ID == artistID {
				if songs.ExistsByID(song.ID) == false {
					songs = append(songs, song)
				}
			}
		}
	}

	return songs, nil
}

func (s *InternalGeniusProvider) getLyrics(songInfo GeniusSongInfo) (string, error) {
	return s.getLyricsFromPath(songInfo.LyricsPath)
}

func (s *InternalGeniusProvider) getLyricsFromPath(lyricsPath string) (string, error) {
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
