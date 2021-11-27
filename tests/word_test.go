package tests

import (
	"github.com/marosiak/WordFinder/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordTrimming(t *testing.T) {
	word := internal.Word("testąćźżĄ")
	word = word.TrimSpecials()

	assert.Equal(t, internal.Word("test"), word)
}

func TestWordOccurancesTrimming(t *testing.T) {
	wordOccurances := internal.WordsOccurrences{
		"testąćźżół": 1,
	}

	wordOccurances = wordOccurances.TrimSpecials()

	assert.Equal(t, 1, wordOccurances["test"])
}
