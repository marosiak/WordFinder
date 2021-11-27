package internal

import (
	"regexp"
	"strings"
	"unicode"
)

type Lyrics string

var MarksChars = []string{",", ".", "!", "?"}
var BannedChars = []string{"()", "[]", "{}"}

func replaceEach(list []string, s string, new string) string {
	for _, v := range list {
		s = strings.ReplaceAll(s, v, new)
	}
	return s
}

func removeParentheses(s string) string {
	return regexp.MustCompile(`\[(.*?)\]|\((.*?)\)`).ReplaceAllString(s, "")
}

func (l Lyrics) Normalised() string {
	lyrics := replaceEach(MarksChars, string(l), "") // TODO: Zmienić MarksChars na []string i wtedy dodać tam 4/MSP
	lyrics = strings.ReplaceAll(lyrics, " ", " ")    // No idea why it appears in lyrics, but I'll leave it there for a sec..

	lyrics = removeParentheses(lyrics)

	previousChar := " "
	output := ""
	for _, char := range lyrics {
		if unicode.IsUpper(char) && previousChar != " " {
			output = output + " " + string(char)
			previousChar = string(char)
			continue
		}
		output = output + string(char)
		previousChar = string(char)
	}

	return output
}

type Word string

func (w Word) TrimSpecials() Word {
	outputWord := w

	for _, bannedChar := range "ĄĆĘŁŃÓŚŹŻ" {
		outputWord = Word(strings.ReplaceAll(string(outputWord), string(bannedChar), ""))
		outputWord = Word(strings.ReplaceAll(string(outputWord), strings.ToLower(string(bannedChar)), ""))
	}

	return outputWord
}

type WordsOccurrences map[Word]int

func (w WordsOccurrences) TrimSpecials() WordsOccurrences {
	outputWordOccurrences := make(WordsOccurrences)

	for word, occurances := range w {
		outputWordOccurrences[word.TrimSpecials()] = occurances
	}
	return outputWordOccurrences
}

func (w WordsOccurrences) ContainsWord(word Word) bool {
	return w[word] > 0
}

func (w WordsOccurrences) ContainsOneOfWords(words []Word) bool {
	wordOccurrences := w.TrimSpecials()
	for _, word := range words {
		if wordOccurrences[word.TrimSpecials()] > 0 {
			return true
		}
	}
	return false
}

func (w WordsOccurrences) Append(theMap WordsOccurrences) WordsOccurrences {
	output := w

	for k, v := range theMap {
		k = Word(replaceEach(BannedChars, string(k), ""))
		occ, ok := output[k]
		if ok == true {
			v = occ + v
		}

		output[k] = v
	}
	return output
}

func splitBySeparators(s string, separators []string) []string {
	for _, separator := range separators {
		s = strings.ReplaceAll(s, separator, " ")
	}

	return strings.Split(s, " ")
}

func (l Lyrics) FindWords() WordsOccurrences {
	lyrics := l.Normalised()

	output := make(WordsOccurrences)
	for _, word := range splitBySeparators(lyrics, []string{" ", "\n"}) {
		if len(word) <= 2 { // Don't count it because it's too short
			continue
		}
		word := Word(strings.ToLower(word))

		occ, ok := output[word]
		if ok == true {
			output[word] = occ + 1
		} else {
			output[word] = 1
		}
	}

	return output
}
