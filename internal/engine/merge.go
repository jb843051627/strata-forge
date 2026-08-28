package engine

import "github.com/jb843051627/strata-forge/internal/model"

func MergeMeasurements(primary, secondary []model.Measurement) []model.Measurement {
	result := primary
	for _, item := range secondary {
		found := false
		for i := range result {
			if result[i].ID == item.ID && item.ID != 0 {
				result[i] = item
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}
	return result
}

func MergeAlerts(primary, secondary []model.Alert) []model.Alert {
	result := model.CloneAlerts(primary)
	for _, item := range secondary {
		found := false
		for i := range result {
			if result[i].ID == item.ID && item.ID != 0 {
				result[i] = item
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}
	return model.SortAlerts(result)
}
