package engine

import "github.com/jb843051627/strata-forge/internal/model"

func BuildSampleGraph(sample model.Sample, layers []model.Layer, measurements []model.Measurement) model.Graph {
	graph := model.NewGraph()
	sampleID := "sample:" + formatID(sample.ID)
	graph.AddNode(model.GraphNode{ID: sampleID, Kind: "sample", State: sample.Status})
	for _, layer := range layers {
		layerID := "layer:" + formatID(layer.ID)
		graph.AddNode(model.GraphNode{ID: layerID, Kind: "layer", State: layer.Status})
		graph.AddEdge(model.GraphEdge{From: sampleID, To: layerID, Kind: "contains"})
		for _, measurement := range measurements {
			if measurement.LayerID != layer.ID {
				continue
			}
			measurementID := "measurement:" + formatID(measurement.ID)
			graph.AddNode(model.GraphNode{ID: measurementID, Kind: measurement.Kind, State: measurement.Status})
			graph.AddEdge(model.GraphEdge{From: layerID, To: measurementID, Kind: "measured"})
		}
	}
	return graph
}

func formatID(value int64) string {
	if value < 0 {
		return "negative"
	}
	return fmtInt(value)
}
