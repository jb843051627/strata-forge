package model

import "fmt"

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c Coordinates) Valid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 && c.Longitude >= -180 && c.Longitude <= 180
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%.5f,%.5f", c.Latitude, c.Longitude)
}

type SiteReference struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
	Elevation   float64     `json:"elevation"`
}

func (s SiteReference) Valid() bool {
	return s.Name != "" && s.Coordinates.Valid()
}
