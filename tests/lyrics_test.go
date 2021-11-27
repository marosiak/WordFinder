package tests

import (
	"github.com/marosiak/WordFinder/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendLyrics(t *testing.T) {
	firstWordsOccMap := internal.WordsOccurrences{
		"word1": 1,
		"word2": 5,
	}

	secondWordsOccMap := internal.WordsOccurrences{
		"word1": 1,
		"word3": 99,
	}

	result := firstWordsOccMap.Append(secondWordsOccMap)

	assert.Equal(t, 2, result["word1"])
	assert.Equal(t, 5, result["word2"])
	assert.Equal(t, 99, result["word3"])
}

func TestFindLyricsWords(t *testing.T) {
	lyrics := internal.Lyrics(
		"aaa bbb\nbbb",
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, 1, wordsMap["aaa"])
	assert.Equal(t, 2, wordsMap["bbb"])
}

func TestFindLyricsWordsSpecial(t *testing.T) {
	lyrics := internal.Lyrics(
		"[test \"test (test)\"] aaa bbb bbb",
	)
	wordsMap := lyrics.FindWords()
	assert.Equal(t, 1, wordsMap["aaa"])
	assert.Equal(t, 2, wordsMap["bbb"])
	assert.Equal(t, 0, wordsMap["test"])
}

func TestFindLyricsWordsSpecialNoSpace(t *testing.T) {
	lyrics := internal.Lyrics(
		"[test \"test (test)\"]aaa bbb bbb",
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, 1, wordsMap["aaa"])
	assert.Equal(t, 2, wordsMap["bbb"])
	assert.Equal(t, 0, wordsMap["test"])
}

func TestFindLyricsWordsSpecialEnter(t *testing.T) {
	lyrics := internal.Lyrics(
		"[test \"test (test)\"]\naaa bbb (test) [test] bbb",
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, 1, wordsMap["aaa"])
	assert.Equal(t, 2, wordsMap["bbb"])
	assert.Equal(t, 0, wordsMap["test"])
}

func TestFindLyricsWordsUpperCaseSplit(t *testing.T) {
	lyrics := internal.Lyrics(
		"test testTest", // there should be 3 "test" words
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, 3, wordsMap["test"])
}
