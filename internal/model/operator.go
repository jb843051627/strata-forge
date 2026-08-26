package model

import "strings"

type Operator struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Specialty []string `json:"specialty"`
}

func (o Operator) Normalized() Operator {
	o.ID = strings.ToLower(strings.TrimSpace(o.ID))
	o.Name = strings.TrimSpace(o.Name)
	for i, specialty := range o.Specialty {
		o.Specialty[i] = strings.ToLower(strings.TrimSpace(specialty))
	}
	return o
}

func (o Operator) CanMeasure(kind string) bool {
	kind = strings.ToLower(strings.TrimSpace(kind))
	for _, specialty := range o.Specialty {
		if specialty == kind {
			return true
		}
	}
	return false
}
