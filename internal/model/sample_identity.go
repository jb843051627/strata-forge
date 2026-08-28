package model

import (
	"fmt"
	"regexp"
	"strings"
)

var sampleCodePattern = regexp.MustCompile(`^[A-Z]{2,8}-[A-Z0-9-]{3,32}$`)

func NormalizeSampleCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func ValidateSampleCode(value string) error {
	value = NormalizeSampleCode(value)
	if !sampleCodePattern.MatchString(value) {
		return fmt.Errorf("%w: invalid sample code", ErrInvalidInput)
	}
	return nil
}

func SampleKey(site, code string) string {
	return strings.ToLower(strings.TrimSpace(site)) + ":" + NormalizeSampleCode(code)
}

func SameSample(a, b Sample) bool {
	return a.ID == b.ID && a.ID != 0
}
