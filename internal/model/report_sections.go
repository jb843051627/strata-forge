package model

type ReportSection struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Order   int    `json:"order"`
	Visible bool   `json:"visible"`
}

func VisibleSections(sections []ReportSection) []ReportSection {
	result := make([]ReportSection, 0, len(sections))
	for _, section := range sections {
		if section.Visible {
			result = append(result, section)
		}
	}
	return result
}

func DefaultSections() []ReportSection {
	return []ReportSection{
		{Title: "样品概览", Order: 1, Visible: true},
		{Title: "分层记录", Order: 2, Visible: true},
		{Title: "质量复核", Order: 3, Visible: true},
		{Title: "年代解释", Order: 4, Visible: true},
	}
}
