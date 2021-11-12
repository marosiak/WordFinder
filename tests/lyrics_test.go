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
		"Lorem ipsum dolor sit.\n" +
			"lorem ipson, polon sit. lorem is",
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, wordsMap["lorem"], 3)
	assert.Equal(t, wordsMap["sit"], 2)
	assert.Equal(t, wordsMap["ipsum"], 1)
}
