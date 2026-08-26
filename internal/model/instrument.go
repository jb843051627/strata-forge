package model

import "strings"

type Instrument struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Serial string `json:"serial"`
	Active bool   `json:"active"`
}

func (i Instrument) Normalize() Instrument {
	i.ID = strings.ToLower(strings.TrimSpace(i.ID))
	i.Kind = strings.ToLower(strings.TrimSpace(i.Kind))
	i.Serial = strings.TrimSpace(i.Serial)
	return i
}

func (i Instrument) Usable() bool {
	return i.ID != "" && i.Kind != "" && i.Serial != "" && i.Active
}
