package engine

import "github.com/jb843051627/strata-forge/internal/model"

func SiteDistance(a, b model.Coordinates) float64 {
	latDelta := a.Latitude - b.Latitude
	longDelta := a.Longitude - b.Longitude
	return (latDelta*latDelta + longDelta*longDelta) * 111.0
}

func InSameRegion(a, b model.Coordinates, radius float64) bool {
	return SiteDistance(a, b) <= radius
}
