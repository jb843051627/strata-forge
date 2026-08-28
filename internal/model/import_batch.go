package model

import "strings"

type ImportRow struct {
	Line        int
	SampleCode  string
	Site        string
	TopDepth    float64
	BottomDepth float64
	Material    string
}

func NormalizeImportRow(row ImportRow) ImportRow {
	row.SampleCode = NormalizeSampleCode(row.SampleCode)
	row.Site = strings.TrimSpace(row.Site)
	row.Material = strings.ToLower(strings.TrimSpace(row.Material))
	return row
}

func ImportRowValid(row ImportRow) bool {
	return row.Line > 0 && row.SampleCode != "" && row.Site != "" && row.BottomDepth > row.TopDepth
}

func GroupImportRows(rows []ImportRow) map[string][]ImportRow {
	groups := make(map[string][]ImportRow)
	for _, row := range rows {
		key := NormalizeSampleCode(row.SampleCode)
		groups[key] = append(groups[key], NormalizeImportRow(row))
	}
	return groups
}
