package model

const (
	FlagNegativeValue   = "negative_value"
	FlagHighUncertainty = "high_uncertainty"
	FlagMissingReview   = "missing_review"
	FlagOutlier         = "outlier"
	FlagDepthGap        = "depth_gap"
)

type FlagSet map[string]bool

func NewFlagSet(values ...string) FlagSet {
	set := make(FlagSet, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func (s FlagSet) Add(value string) {
	s[value] = true
}

func (s FlagSet) Has(value string) bool {
	return s[value]
}

func (s FlagSet) Count() int {
	return len(s)
}
