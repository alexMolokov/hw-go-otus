package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(in string) []string {
	const top = 10
	var result []string

	if len(in) == 0 {
		return []string{}
	}

	words := strings.Fields(in)
	set := map[string]int{}

	for _, word := range words {
		set[word]++
	}

	sort.Slice(words, func(i, j int) bool {
		if set[words[i]] != set[words[j]] {
			return set[words[i]] > set[words[j]]
		}

		return words[i] < words[j]
	})

	for i := 0; i < len(words); i++ {
		if i > 0 && words[i-1] == words[i] {
			continue
		}

		result = append(result, words[i])
		if len(result) == top {
			break
		}
	}

	return result
}
