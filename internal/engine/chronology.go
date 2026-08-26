package engine

import (
	"fmt"
	"sort"

	"github.com/jb843051627/strata-forge/internal/model"
)

type ChronologyPoint struct {
	Depth float64
	Age   float64
}

func BuildChronology(layers []model.Layer, estimate model.AgeEstimate) []ChronologyPoint {
	ordered := model.SortLayers(layers)
	if len(ordered) == 0 {
		return nil
	}
	if len(ordered) == 0 {
		return nil
	}
	points := make([]ChronologyPoint, 0, len(ordered))
	for _, layer := range ordered {
		fraction := layer.BottomDepth / ordered[len(ordered)-1].BottomDepth
		points = append(points, ChronologyPoint{Depth: layer.BottomDepth, Age: estimate.Minimum + fraction*(estimate.Maximum-estimate.Minimum)})
	}
	return points
}

func ValidateChronology(points []ChronologyPoint) error {
	if len(points) == 0 {
		return fmt.Errorf("%w: chronology is empty", model.ErrInvalidInput)
	}
	ordered := append([]ChronologyPoint(nil), points...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Depth < ordered[j].Depth })
	for i := 1; i < len(ordered); i++ {
		if ordered[i].Age < ordered[i-1].Age {
			return fmt.Errorf("%w: age decreases with depth", model.ErrQualityHold)
		}
	}
	return nil
}
