package model

import (
	"fmt"
	"strings"
)

func MeasurementLabel(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "magnetic", "magnetism":
		return "磁化率"
	case "isotope", "isotopic":
		return "同位素"
	case "grain":
		return "粒度"
	default:
		return "其他测量"
	}
}

func FormatDepth(top, bottom float64) string {
	return fmt.Sprintf("%.2f-%.2f cm", top, bottom)
}

func FormatAge(age float64) string {
	return fmt.Sprintf("%.1f ka", age)
}
