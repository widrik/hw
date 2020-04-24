package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
)

const top = 10

type wordStatistic struct {
	Word  string
	Count int
}

var regexpSplit = regexp.MustCompile(`[\s\t\r\n]+`)

func splitTextToParts(inputText string) []string {
	return regexpSplit.Split(inputText, -1)
}

func Top10(inputText string) []string {
	if len(inputText) == 0 {
		return nil
	}

	parts := splitTextToParts(inputText)

	wordsList := make(map[string]int)
	for _, part := range parts {
		if len(part) > 0 {
			wordsList[part]++
		}
	}

	var words = make([]wordStatistic, 0, len(wordsList))
	for word, count := range wordsList {
		words = append(words, wordStatistic{word, count})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	var limit int
	if len(words) < top {
		limit = len(words)
	} else {
		limit = top
	}

	var result = make([]string, 0, limit)
	for _, word := range words[:limit] {
		result = append(result, word.Word)
	}

	return result
}
