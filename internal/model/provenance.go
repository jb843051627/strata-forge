package model

import (
	"strings"
	"time"
)

type Provenance struct {
	Source     string    `json:"source"`
	Instrument string    `json:"instrument"`
	OperatorID string    `json:"operator_id"`
	CapturedAt time.Time `json:"captured_at"`
	Checksum   string    `json:"checksum"`
}

func (p Provenance) Complete() bool {
	return strings.TrimSpace(p.Source) != "" && strings.TrimSpace(p.Instrument) != "" && !p.CapturedAt.IsZero()
}

func (p Provenance) Normalize() Provenance {
	p.Source = strings.TrimSpace(p.Source)
	p.Instrument = strings.TrimSpace(p.Instrument)
	p.OperatorID = strings.ToLower(strings.TrimSpace(p.OperatorID))
	p.Checksum = strings.TrimSpace(p.Checksum)
	return p
}
