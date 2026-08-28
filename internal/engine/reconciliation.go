package engine

import (
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func ReconcileLayerDepth(sample model.Sample, layers []model.Layer) error {
	if len(layers) == 0 {
		return fmt.Errorf("%w: no layers", model.ErrInvalidInput)
	}
	ordered := model.SortLayers(layers)
	previous := sample.DepthStart
	for _, layer := range ordered {
		if layer.TopDepth < previous || layer.BottomDepth > sample.DepthEnd {
			return fmt.Errorf("%w: layer %d outside depth sequence", model.ErrQualityHold, layer.ID)
		}
		previous = layer.BottomDepth
	}
	return nil
}

func ReconcileMeasurementLayer(measurement model.Measurement, layer model.Layer) error {
	if measurement.LayerID != layer.ID {
		return fmt.Errorf("%w: measurement layer mismatch", model.ErrInvalidInput)
	}
	return nil
}
