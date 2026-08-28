package engine

import "github.com/jb843051627/strata-forge/internal/model"

func SnapshotMeasurements(in []model.Measurement) []model.Measurement {
	return model.CloneMeasurements(in)
}

func SnapshotAlerts(in []model.Alert) []model.Alert {
	return model.CloneAlerts(in)
}

func SnapshotLayers(in []model.Layer) []model.Layer {
	return model.CloneLayers(in)
}
