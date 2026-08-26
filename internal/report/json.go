package report

import (
	"encoding/json"
	"strconv"

	"github.com/jb843051627/strata-forge/internal/model"
)

type Envelope struct {
	Kind    string `json:"kind"`
	Version string `json:"version"`
	Data    any    `json:"data"`
}

func EncodeEnvelope(kind, version string, data any) ([]byte, error) {
	return json.MarshalIndent(Envelope{Kind: kind, Version: version, Data: data}, "", "  ")
}

func EncodeReport(report model.Report) ([]byte, error) {
	return EncodeEnvelope("strata-forge-report", fmtVersion(report.Version), report)
}

func fmtVersion(version int) string {
	if version <= 0 {
		return "draft"
	}
	return "v" + strconv.Itoa(version)
}
