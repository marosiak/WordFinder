package internal

import (
	"errors"
	"github.com/marosiak/WordFinder/config"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SongInfo struct {
	AuthorName     string
	Title          string
	LyricsEndpoint string
	GeniusID       int
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
	GetSongFromInfo(songInfo SongInfo) (Song, error)

	GetSongsInfosByArtist(artistName string) ([]SongInfo, error)
	GetSongsByArtist(artistName string) ([]Song, error)
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

func (s *InternalLyricsService) GetSongInfoByName(songName string) (SongInfo, error) {
	searchResults, err := s.geniusProvider.Search(songName)
	if err != nil {
		return SongInfo{}, err
	}

	song := searchResults[0]

	return SongInfo{
		GeniusID:       song.ID,
		AuthorName:     song.PrimaryArtist.Name,
		Title:          song.FullTitle,
		LyricsEndpoint: song.LyricsEndpoint,
	}, nil
}

func (s *InternalLyricsService) GetSongInfoByID(id int) (SongInfo, error) {
	geniusSong, err := s.geniusProvider.GetSongInfoByID(id)
	if err != nil {
		return SongInfo{}, err
	}

	return SongInfo{
		GeniusID:       geniusSong.ID,
		AuthorName:     geniusSong.PrimaryArtist.Name,
		Title:          geniusSong.FullTitle,
		LyricsEndpoint: geniusSong.LyricsPath,
	}, nil
}

func (s *InternalLyricsService) GetSongByName(songName string) (Song, error) {
	geniusSong, err := s.geniusProvider.GetSongByName(songName)
	if err != nil {
		return Song{}, err
	}

	return Song{
		Lyrics: geniusSong.Lyrics,
		Info: SongInfo{
			GeniusID:       geniusSong.Info.ID,
			AuthorName:     geniusSong.Info.PrimaryArtist.Name,
			Title:          geniusSong.Info.FullTitle,
			LyricsEndpoint: geniusSong.Info.LyricsPath,
		},
	}, nil
}

func (s *InternalLyricsService) GetSongFromInfo(songInfo SongInfo) (Song, error) {
	song, err := s.geniusProvider.GetSongByID(songInfo.GeniusID)
	if err != nil {
		return Song{}, err
	}

	return Song{
		Info:   songInfo,
		Lyrics: song.Lyrics,
	}, nil
}

func (s *InternalLyricsService) GetSongsInfosByArtist(artistName string) ([]SongInfo, error) {
	primaryArtistID, err := s.findArtistID(artistName)
	if err != nil {
		s.logger.WithError(err).Error("GetSongsInfosByArtist find artist ID for ", artistName)
	}

	foundSongs, err := s.geniusProvider.FindSongInfosByArtistID(primaryArtistID)
	if err != nil {
		return []SongInfo{}, err
	}

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

func (s *InternalLyricsService) GetSongsFromSongInfos(songInfos []SongInfo) ([]Song, error) {
	var ids []int
	for _, songInfo := range songInfos {
		ids = append(ids, songInfo.GeniusID)
	}

	geniusSongs, err := s.geniusProvider.GetSongsByIDs(ids)
	if err != nil {
		return []Song{}, nil
	}

	var songs []Song
	for _, geniusSong := range geniusSongs {
		songs = append(songs, Song{
			Info: SongInfo{
				AuthorName:     geniusSong.Info.PrimaryArtist.Name,
				Title:          geniusSong.Info.FullTitle,
				LyricsEndpoint: geniusSong.Info.LyricsPath,
				GeniusID:       geniusSong.Info.ID,
			},
			Lyrics: geniusSong.Lyrics,
		})
	}
	return songs, nil
}

func (s *InternalLyricsService) GetSongsByArtist(artistName string) ([]Song, error) {
	artistID, err := s.findArtistID(artistName)
	if err != nil {
		return []Song{}, err
	}

	geniusSongs, err := s.geniusProvider.FindSongsByArtistID(artistID)
	if err != nil {
		return nil, err
	}

	var songs []Song
	for _, geniusSong := range geniusSongs {
		songs = append(songs, Song{
			Info: SongInfo{
				AuthorName:     geniusSong.Info.PrimaryArtist.Name,
				Title:          geniusSong.Info.FullTitle,
				LyricsEndpoint: geniusSong.Info.LyricsPath,
			},
			Lyrics: geniusSong.Lyrics,
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
