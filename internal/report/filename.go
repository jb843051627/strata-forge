package report

import (
	"fmt"
	"strings"
)

func SafeFilename(code string, version int) string {
	clean := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, strings.TrimSpace(code))
	if clean == "" {
		clean = "sample"
	}
	return fmt.Sprintf("%s-v%d.md", clean, version)
}
