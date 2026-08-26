package engine

type RuleSet struct {
	MinimumScore   float64
	MaxUncertainty float64
	MaxGap         float64
	MinimumDepth   float64
}

func DefaultRules() RuleSet {
	return RuleSet{MinimumScore: 70, MaxUncertainty: 0.25, MaxGap: 5, MinimumDepth: 0.01}
}

func (r RuleSet) WithUncertainty(limit float64) RuleSet {
	r.MaxUncertainty = limit
	return r
}

func (r RuleSet) WithMinimumScore(score float64) RuleSet {
	r.MinimumScore = score
	return r
}

func (r RuleSet) Valid() bool {
	return r.MinimumScore > 0 && r.MinimumScore <= 100 && r.MaxUncertainty > 0 && r.MaxGap > 0 && r.MinimumDepth > 0
}
