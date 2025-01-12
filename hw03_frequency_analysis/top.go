package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	words := strings.Fields(str)

	wordFrequency := calculateFrequency(words)

	// Creating a slice of keys
	keys := make([]string, 0, len(wordFrequency))
	for word := range wordFrequency {
		keys = append(keys, word)
	}

	// Sorting the keys by frequency
	sort.Slice(keys, func(i, j int) bool {
		if wordFrequency[keys[i]] == wordFrequency[keys[j]] {
			return keys[i] < keys[j]
		}
		return wordFrequency[keys[i]] > wordFrequency[keys[j]]
	})

	// Limiting the result to 10 words
	limit := 10
	if len(keys) < limit {
		limit = len(keys)
	}

	return keys[:limit]
}

func calculateFrequency(words []string) map[string]int {
	wordFrequency := make(map[string]int)
	for _, word := range words {
		wordFrequency[word]++
	}
	return wordFrequency
}
