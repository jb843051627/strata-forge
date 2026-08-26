package model

import "strings"

func ReviewDecisionLabel(decision string) string {
	switch decision {
	case ReviewPass:
		return "通过"
	case ReviewHold:
		return "待补测"
	case ReviewFail:
		return "退回"
	default:
		return strings.ToUpper(decision)
	}
}

func ReviewAllowsReport(decision string) bool {
	return decision == ReviewPass
}

func ReviewNeedsFollowup(decision string) bool {
	return decision == ReviewHold || decision == ReviewFail
}
