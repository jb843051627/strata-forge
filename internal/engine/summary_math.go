package engine

import (
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

func StatusCounts(items []model.Measurement) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return counts
}

func KindCounts(items []model.Measurement) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Kind]++
	}
	return counts
}

func SortedKinds(counts map[string]int) []string {
	result := make([]string, 0, len(counts))
	for key := range counts {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
