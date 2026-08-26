package report

import "github.com/jb843051627/strata-forge/internal/model"

func SummaryTitle(summary *model.SampleSummary) string {
	if summary == nil {
		return "未选择样品"
	}
	return summary.Sample.Code + " 年代解释"
}
