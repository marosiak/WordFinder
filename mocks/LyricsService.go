// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	internal "github.com/marosiak/WordFinder/internal"
	mock "github.com/stretchr/testify/mock"
)

// LyricsService is an autogenerated mock type for the LyricsService type
type LyricsService struct {
	mock.Mock
}

// GetArtist provides a mock function with given fields: artistName
func (_m *LyricsService) GetArtist(artistName string) (internal.Artist, error) {
	ret := _m.Called(artistName)

	var r0 internal.Artist
	if rf, ok := ret.Get(0).(func(string) internal.Artist); ok {
		r0 = rf(artistName)
	} else {
		r0 = ret.Get(0).(internal.Artist)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(artistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongByName provides a mock function with given fields: name
func (_m *LyricsService) GetSongByName(name string) (internal.Song, error) {
	ret := _m.Called(name)

	var r0 internal.Song
	if rf, ok := ret.Get(0).(func(string) internal.Song); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(internal.Song)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongFromInfo provides a mock function with given fields: songInfo
func (_m *LyricsService) GetSongFromInfo(songInfo internal.SongInfo) (internal.Song, error) {
	ret := _m.Called(songInfo)

	var r0 internal.Song
	if rf, ok := ret.Get(0).(func(internal.SongInfo) internal.Song); ok {
		r0 = rf(songInfo)
	} else {
		r0 = ret.Get(0).(internal.Song)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(internal.SongInfo) error); ok {
		r1 = rf(songInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongInfoByID provides a mock function with given fields: id
func (_m *LyricsService) GetSongInfoByID(id int) (internal.SongInfo, error) {
	ret := _m.Called(id)

	var r0 internal.SongInfo
	if rf, ok := ret.Get(0).(func(int) internal.SongInfo); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(internal.SongInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongInfoByName provides a mock function with given fields: name
func (_m *LyricsService) GetSongInfoByName(name string) (internal.SongInfo, error) {
	ret := _m.Called(name)

	var r0 internal.SongInfo
	if rf, ok := ret.Get(0).(func(string) internal.SongInfo); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(internal.SongInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongsByArtist provides a mock function with given fields: artistName
func (_m *LyricsService) GetSongsByArtist(artistName string) ([]internal.Song, error) {
	ret := _m.Called(artistName)

	var r0 []internal.Song
	if rf, ok := ret.Get(0).(func(string) []internal.Song); ok {
		r0 = rf(artistName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(artistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongsFromSongInfos provides a mock function with given fields: songInfos
func (_m *LyricsService) GetSongsFromSongInfos(songInfos []internal.SongInfo) ([]internal.Song, error) {
	ret := _m.Called(songInfos)

	var r0 []internal.Song
	if rf, ok := ret.Get(0).(func([]internal.SongInfo) []internal.Song); ok {
		r0 = rf(songInfos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Song)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]internal.SongInfo) error); ok {
		r1 = rf(songInfos)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongsInfosByArtist provides a mock function with given fields: artistName
func (_m *LyricsService) GetSongsInfosByArtist(artistName string) ([]internal.SongInfo, error) {
	ret := _m.Called(artistName)

	var r0 []internal.SongInfo
	if rf, ok := ret.Get(0).(func(string) []internal.SongInfo); ok {
		r0 = rf(artistName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.SongInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(artistName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
