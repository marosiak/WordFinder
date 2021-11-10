package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendLyrics(t *testing.T) {
	firstWordsOccMap := WordsOccurrences{
		"word1": 1,
		"word2": 5,
	}

	secondWordsOccMap := WordsOccurrences{
		"word1": 1,
		"word3": 99,
	}

	result := firstWordsOccMap.Append(secondWordsOccMap)

	assert.Equal(t, 2, result["word1"])
	assert.Equal(t, 5, result["word2"])
	assert.Equal(t, 99, result["word3"])
}

func TestFindLyricsWords(t *testing.T) {
	lyrics := Lyrics(
		"Lorem ipsum dolor sit.\n" +
			"lorem ipson, polon sit. lorem",
	)

	wordsMap := lyrics.FindWords()
	assert.Equal(t, wordsMap["lorem"], 3)
	assert.Equal(t, wordsMap["sit"], 2)
	assert.Equal(t, wordsMap["ipsum"], 1)
}
