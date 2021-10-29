package internal

import (
	"errors"
	"fmt"
	"github.com/marosiak/WordFinder/config"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type SongInfo struct {
	AuthorName     string
	Title          string
	LyricsEndpoint string
}

type Song struct {
	Info   SongInfo
	Lyrics string
}

func (s *SongInfo) CountOccurances(word string) int {
	return strings.Count(s.Title, word)
}

type LyricsService interface {
	GetSongInfoByName(name string) (SongInfo, error)
	GetSongInfoByID(id int) (SongInfo, error)
	GetSongByName(name string) (Song, error)
	GetAllSongsInfoByArtist(artistName string) ([]SongInfo, error)
	GetSongFromInfo(songInfo SongInfo) (Song, error)
	GetSongsFromSongInfos(songInfos []SongInfo) ([]Song, error)
}

type InternalLyricsService struct {
	geniusProvider GeniusProvider
	logger         *log.Entry
	cfg            *config.Config
}

func NewLyricsService(cfg *config.Config, geniusProvider GeniusProvider, logger *log.Entry) *InternalLyricsService {
	return &InternalLyricsService{geniusProvider: geniusProvider, logger: logger, cfg: cfg}
}

var _ LyricsService = &InternalLyricsService{}

func (s *InternalLyricsService) GetSongInfoByID(id int) (SongInfo, error) {
	geniusSong, err := s.geniusProvider.GetSongByID(id)
	if err != nil {
		return SongInfo{}, err
	}

	return SongInfo{
		AuthorName:     geniusSong.PrimaryArtist.Name,
		Title:          geniusSong.FullTitle,
		LyricsEndpoint: geniusSong.LyricsPath,
	}, nil
}

func (s *InternalLyricsService) GetSongInfoByName(songName string) (SongInfo, error) {
	searchResults, err := s.geniusProvider.Search(songName)
	if err != nil {
		return SongInfo{}, err
	}

	song := searchResults[0]

	return SongInfo{
		AuthorName:     song.PrimaryArtist.Name,
		Title:          song.FullTitle,
		LyricsEndpoint: song.LyricsEndpoint,
	}, nil
}

func (s *InternalLyricsService) GetSongByName(songName string) (Song, error) {
	songInfo, err := s.GetSongInfoByName(songName)
	if err != nil {
		return Song{}, err
	}

	lyrics, err := s.geniusProvider.GetLyricsFromPath(songInfo.LyricsEndpoint)
	if err != nil {
		return Song{}, err
	}

	return Song{
		Info:   songInfo,
		Lyrics: lyrics,
	}, nil
}

func (s *InternalLyricsService) GetSongsFromSongInfos(songInfos []SongInfo) ([]Song, error) {
	songCh := make(chan Song)
	wg := sync.WaitGroup{}

	for _, songInfo := range songInfos {
		wg.Add(1)
		songInfo := songInfo
		go func() {
			defer wg.Done()
			song, err := s.GetSongFromInfo(songInfo)
			if err != nil {
				s.logger.WithError(err)
				return
			}
			songCh <- song
		}()
	}

	go func() {
		wg.Wait()
		close(songCh)
	}()

	var songs []Song
	for song := range songCh {
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *InternalLyricsService) GetSongFromInfo(songInfo SongInfo) (Song, error) {
	lyrics, err := s.geniusProvider.GetLyricsFromPath(songInfo.LyricsEndpoint)
	if err != nil {
		return Song{}, err
	}

	return Song{
		Info:   songInfo,
		Lyrics: lyrics,
	}, nil
}

func (s *InternalLyricsService) GetAllSongsByArtist(artistName string) ([]Song, error) {
	primaryArtistID, err := s.findArtistID(artistName)
	if err != nil {
		s.logger.WithError(err).Error("GetAllSongsByArtist cannot find artist ID for ", artistName)
	}

	geniusSongs, err := s.geniusProvider.FindSongsByArtistID(primaryArtistID)
	if err != nil {
		return nil, err
	}

	var songs []Song

	for _, geniusSong := range geniusSongs {
		song, err := s.geniusProvider.GetSongByID(geniusSong.ID)
		if err != nil {
			s.logger.WithError(err)
			continue
		}

		lyrics, err := s.geniusProvider.GetLyrics(song)
		if err != nil {
			s.logger.WithError(err).Error("GetAllSongsByArtist cannot GetLyrics")
		}

		songs = append(songs, Song{
			Info: SongInfo{
				AuthorName:     song.PrimaryArtist.Name,
				Title:          song.FullTitle,
				LyricsEndpoint: song.LyricsPath,
			},
			Lyrics: lyrics,
		})
	}

	return songs, nil
}

func (s *InternalLyricsService) findArtistID(desiredArtistName string) (int, error) {
	searchResults, err := s.geniusProvider.Search(desiredArtistName)
	if err != nil {
		return 0, err
	}

	desiredArtistName = strings.ToLower(desiredArtistName)
	for _, result := range searchResults {
		primaryArtistName := strings.ToLower(result.PrimaryArtist.Name)

		if strings.Contains(primaryArtistName, desiredArtistName) {
			return result.PrimaryArtist.ID, nil
		}
	}
	return 0, errors.New("artist not found")
}

func (s *InternalLyricsService) GetAllSongsInfoByArtist(artistName string) ([]SongInfo, error) {
	primaryArtistID, err := s.findArtistID(artistName)
	if err != nil {
		s.logger.WithError(err).Error("GetAllSongsInfoByArtist find artist ID for ", artistName)
	}

	foundSongs, err := s.geniusProvider.FindSongsByArtistID(primaryArtistID)
	if err != nil {
		return []SongInfo{}, err
	}
	fmt.Printf("Found: %d songs\n", len(foundSongs))

	var songs []SongInfo
	for _, song := range foundSongs {
		songs = append(songs, SongInfo{
			AuthorName:     song.PrimaryArtist.Name,
			Title:          song.FullTitle,
			LyricsEndpoint: song.LyricsPath,
		})
	}

	return songs, nil
}
