package engine

import (
	"fmt"
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

func ValidateImportRows(rows []model.ImportRow) []string {
	issues := make([]string, 0)
	for _, row := range rows {
		if !model.ImportRowValid(row) {
			issues = append(issues, fmt.Sprintf("line %d invalid", row.Line))
		}
	}
	return issues
}

func SortImportRows(rows []model.ImportRow) []model.ImportRow {
	result := append([]model.ImportRow(nil), rows...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].SampleCode == result[j].SampleCode {
			return result[i].TopDepth < result[j].TopDepth
		}
		return result[i].SampleCode < result[j].SampleCode
	})
	return result
}

func BuildIntervals(rows []model.ImportRow) []model.DepthInterval {
	result := make([]model.DepthInterval, 0, len(rows))
	for _, row := range rows {
		result = append(result, model.DepthInterval{Top: row.TopDepth, Bottom: row.BottomDepth})
	}
	return result
}
