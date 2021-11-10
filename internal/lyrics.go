package internal

import (
	"sort"
	"strings"
)

type Lyrics string

var SpecialChars = ",.!?"

func (l Lyrics) Raw() string {
	lyrics := string(l)
	for _, specialChar := range SpecialChars {
		lyrics = strings.ReplaceAll(lyrics, string(specialChar), "")
	}
	lyrics = strings.ReplaceAll(lyrics, "\n", " ")
	return lyrics
}

type WordsOccurrences map[string]int

type KeyValue struct {
	Key   string
	Value int
}

func (w WordsOccurrences) SortedKeyValueList() []KeyValue {
	var ss []KeyValue
	for k, v := range w {
		ss = append(ss, KeyValue{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	return ss
}

func (w WordsOccurrences) Append(theMap map[string]int) WordsOccurrences {
	output := w
	for k, v := range theMap {
		occ, ok := output[k]
		if ok == true {
			v = occ + v
			println(v)
		}

		output[k] = v
	}
	return output
}

func (l Lyrics) FindWords() WordsOccurrences {
	lyrics := l.Raw()
	splitted := append(strings.Split(lyrics, " "))

	output := make(map[string]int)
	for _, word := range splitted {
		if len(word) <= 2 {
			continue
		}
		word := strings.ToLower(word)
		occ, ok := output[word]
		if ok == true {
			output[word] = occ + 1
		} else {
			output[word] = 1
		}
	}

	return output
}
