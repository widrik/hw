package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
)

const top = 10

type WordStatistic struct {
	Word  string
	Count int
}

func splitTextToParts(inputText string) []string {
	return regexp.MustCompile(`[\s\t\r\n]+`).Split(inputText, -1)
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

	words := []WordStatistic{}
	for word, count := range wordsList {
		words = append(words, WordStatistic{word, count})
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

	result := []string{}
	for _, word := range words[:limit] {
		result = append(result, word.Word)
	}

	return result
}
