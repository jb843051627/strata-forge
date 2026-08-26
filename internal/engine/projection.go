package engine

import (
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

type LayerProjection struct {
	LayerID int64
	Top     float64
	Bottom  float64
	Age     float64
}

func ProjectAges(layers []model.Layer, estimate model.AgeEstimate) []LayerProjection {
	ordered := model.SortLayers(layers)
	result := make([]LayerProjection, 0, len(ordered))
	for _, layer := range ordered {
		fraction := 0.0
		if estimate.Maximum != estimate.Minimum {
			fraction = layer.BottomDepth / (ordered[len(ordered)-1].BottomDepth + 0.0001)
		}
		result = append(result, LayerProjection{LayerID: layer.ID, Top: layer.TopDepth, Bottom: layer.BottomDepth, Age: estimate.Minimum + fraction*(estimate.Maximum-estimate.Minimum)})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Bottom < result[j].Bottom })
	return result
}

func ProjectionRange(values []LayerProjection) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	minimum, maximum := values[0].Age, values[0].Age
	for _, value := range values[1:] {
		if value.Age < minimum {
			minimum = value.Age
		}
		if value.Age > maximum {
			maximum = value.Age
		}
	}
	return minimum, maximum
}
