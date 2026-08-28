package model

import "sort"

func SortLayers(items []Layer) []Layer {
	result := CloneLayers(items)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].TopDepth == result[j].TopDepth {
			return result[i].Sequence < result[j].Sequence
		}
		return result[i].TopDepth < result[j].TopDepth
	})
	return result
}

func SortMeasurements(items []Measurement) []Measurement {
	result := CloneMeasurements(items)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].LayerID == result[j].LayerID {
			return result[i].ID < result[j].ID
		}
		return result[i].LayerID < result[j].LayerID
	})
	return result
}

func SortAlerts(items []Alert) []Alert {
	result := CloneAlerts(items)
	priority := map[string]int{"critical": 0, "warning": 1, "info": 2}
	sort.SliceStable(result, func(i, j int) bool {
		return priority[result[i].Severity] < priority[result[j].Severity]
	})
	return result
}
