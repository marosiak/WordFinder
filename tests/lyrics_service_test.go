package tests

import (
	"errors"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
)

var anyError = errors.New("error")

func getLyricsServiceAndGeniusProvider() (*mocks.GeniusProvider, *internal.InternalLyricsService) {
	cfg := GetConfig()

	geniusProvider := &mocks.GeniusProvider{}
	return geniusProvider, internal.NewLyricsService(cfg, geniusProvider, log.NewEntry(log.New()))
}

func TestGetArtistSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetArtist", "the_artist").Return(
		internal.GeniusArtist{
			ID:      1,
			ApiPath: "/1/",
			Name:    "the_artist_full_name",
		}, nil,
	)

	artist, err := lyricsService.GetArtist("the_artist")
	assert.NoError(t, err)

	assert.Equal(t, 1, artist.GeniusID)
	assert.Equal(t, "the_artist_full_name", artist.Name)
}

func TestGetArtistError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetArtist", "the_artist").Return(
		internal.GeniusArtist{}, anyError,
	)

	_, err := lyricsService.GetArtist("the_artist")
	assert.Error(t, err, anyError)
}

func TestGetSongFromInfoSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongByID", 1).Return(
		internal.GeniusSong{
			Lyrics: "the lyrics are here",
			Info: internal.GeniusSongInfo{
				ID:         1,
				LyricsPath: "path",
				FullTitle:  "title",
				PrimaryArtist: internal.GeniusArtist{
					ID:      2,
					ApiPath: "path2",
					Name:    "artist",
				},
				LyricsState: "state",
			},
		}, nil,
	)

	artist, err := lyricsService.GetSongFromInfo(internal.SongInfo{
		AuthorName:     "artist",
		Title:          "title",
		LyricsEndpoint: "path",
		GeniusID:       1,
	})

	assert.NoError(t, err)
	assert.Equal(t, artist.Info.Title, "title")
}

func TestGetSongFromInfoError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongByID", 1).Return(
		internal.GeniusSong{}, anyError,
	)

	_, err := lyricsService.GetSongFromInfo(internal.SongInfo{
		AuthorName:     "artist",
		Title:          "title",
		LyricsEndpoint: "path",
		GeniusID:       1,
	})

	assert.Error(t, err, anyError)
}

func TestGetSongInfoByNameSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("Search", "the_song").Return(
		[]internal.GeniusSearchResult{
			{
				ID:             1,
				ApiPath:        "a_path",
				FullTitle:      "the_song",
				LyricsEndpoint: "/lyrics/",
				PrimaryArtist: internal.GeniusArtist{
					ID:      1,
					ApiPath: "1",
					Name:    "1",
				},
			},
		}, nil,
	)

	song, err := lyricsService.GetSongInfoByName("the_song")
	assert.NoError(t, err)
	assert.Equal(t, "the_song", song.Title)
}

func TestGetSongInfoByNameError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("Search", "the_song").Return(
		[]internal.GeniusSearchResult{}, anyError,
	)

	_, err := lyricsService.GetSongInfoByName("the_song")
	assert.Error(t, err, anyError)
}

func TestGetSongByNameSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongByName", "the_song").Return(
		internal.GeniusSong{
			Lyrics: "the lyrics",
			Info: internal.GeniusSongInfo{
				ID:        1,
				FullTitle: "the_song",
			},
		}, nil,
	)

	song, err := lyricsService.GetSongByName("the_song")
	assert.NoError(t, err)
	assert.Equal(t, "the lyrics", song.Info.Title)
}
