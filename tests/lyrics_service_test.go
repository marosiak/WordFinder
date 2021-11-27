package tests

import (
	"errors"
	"github.com/marosiak/WordFinder/internal"
	"github.com/marosiak/WordFinder/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
				ID:        1,
				PagePath:  "path",
				FullTitle: "title",
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
		AuthorName:   "artist",
		Title:        "title",
		PageEndpoint: "path",
		GeniusID:     1,
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
		AuthorName:   "artist",
		Title:        "title",
		PageEndpoint: "path",
		GeniusID:     1,
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
	assert.Equal(t, "the_song", song.Info.Title)
	assert.Equal(t, internal.Lyrics("the lyrics"), song.Lyrics)
}

func TestGetSongByNameError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongByName", "the_song").Return(
		internal.GeniusSong{}, anyError,
	)

	_, err := lyricsService.GetSongByName("the_song")
	assert.Error(t, anyError, err)
}

func TestGetSongInfoByIDSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongInfoByID", 1).Return(
		internal.GeniusSongInfo{
			ID:            1,
			PagePath:      "lp",
			FullTitle:     "full_tittle",
			PrimaryArtist: internal.GeniusArtist{},
			LyricsState:   "ls",
		}, nil,
	)

	song, err := lyricsService.GetSongInfoByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "full_tittle", song.Title)
	assert.Equal(t, 1, song.GeniusID)
}

func TestGetSongInfoByIDError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetSongInfoByID", mock.Anything).Return(
		internal.GeniusSongInfo{}, anyError,
	)

	_, err := lyricsService.GetSongInfoByID(1)
	assert.Error(t, anyError, err)
}

func TestGetSongsInfosByArtistSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetArtist", "artist").Return(
		internal.GeniusArtist{
			ID:      1,
			ApiPath: "/1/",
			Name:    "the_artist_full_name",
		}, nil,
	)

	geniusProvider.On("GetSongInfosByArtistID", mock.Anything).Return(
		[]internal.GeniusSongInfo{
			{
				ID:        1,
				FullTitle: "song title",
				PrimaryArtist: internal.GeniusArtist{
					ID:   1,
					Name: "artist",
				},
			},
		}, nil,
	)

	song, err := lyricsService.GetSongsInfosByArtist("artist")
	assert.NoError(t, err)
	assert.Equal(t, "song title", song[0].Title)
}

func TestGetSongsInfosByArtistNoArtistError(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetArtist", "artist").Return(
		internal.GeniusArtist{}, anyError,
	)

	_, err := lyricsService.GetSongsInfosByArtist("artist")
	assert.Error(t, err, anyError)
}

func TestGetSongsInfosByArtistErrorGettingArtistSongs(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()

	geniusProvider.On("GetArtist", "artist").Return(
		internal.GeniusArtist{
			ID:      1,
			ApiPath: "/1/",
			Name:    "the_artist_full_name",
		}, nil,
	)

	geniusProvider.On("GetSongInfosByArtistID", mock.Anything).Return(
		[]internal.GeniusSongInfo{}, anyError,
	)

	_, err := lyricsService.GetSongsInfosByArtist("artist")
	assert.Error(t, anyError, err)
}

func TestGetSongsFromSongInfosSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()
	geniusProvider.On("GetSongsByIDs", mock.Anything).Return(
		[]internal.GeniusSong{
			{
				Info: internal.GeniusSongInfo{
					ID:        1,
					FullTitle: "the_title",
				},
			},
			{
				Info: internal.GeniusSongInfo{
					ID:        1,
					FullTitle: "the_title 2",
				},
			},
		}, nil,
	)

	geniusProvider.On("GetSongInfosByArtistID", mock.Anything).Return(
		[]internal.GeniusSongInfo{}, anyError,
	)

	songs, err := lyricsService.GetSongsFromSongInfos([]internal.SongInfo{
		{GeniusID: 1, Title: "the_title"},
		{GeniusID: 2, Title: "the_title 2"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(songs))
}

//func TestGetSongsFromSongInfosError(t *testing.T) {
//	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()
//	geniusProvider.On("GetSongsByIDs", mock.Anything).Return(
//		[]internal.GeniusSong{}, anyError,
//	)
//
//	_, err := lyricsService.GetSongsFromSongInfos([]internal.SongInfo{})
//	assert.Error(t, err, anyError)
//}

func TestGetSongsByArtistSuccess(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()
	geniusProvider.On("GetArtist", mock.Anything).Return(internal.GeniusArtist{
		ID:   2,
		Name: "artist",
	}, nil)

	geniusProvider.On("GetSongsByArtistID", 2).Return([]internal.GeniusSong{
		{}, {},
	}, nil)

	songs, err := lyricsService.GetSongsByArtist("artist")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(songs))
}

func TestGetSongsByArtistErrorGettingArtist(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()
	geniusProvider.On("GetArtist", mock.Anything).Return(internal.GeniusArtist{}, anyError)

	_, err := lyricsService.GetSongsByArtist("artist")
	assert.Error(t, err, anyError)
}

func TestGetSongsByArtistErrorGettingSongs(t *testing.T) {
	geniusProvider, lyricsService := getLyricsServiceAndGeniusProvider()
	geniusProvider.On("GetArtist", mock.Anything).Return(internal.GeniusArtist{
		ID:   2,
		Name: "artist",
	}, nil)

	geniusProvider.On("GetSongsByArtistID", 2).Return([]internal.GeniusSong{}, anyError)

	_, err := lyricsService.GetSongsByArtist("artist")
	assert.Error(t, err, anyError)
}
