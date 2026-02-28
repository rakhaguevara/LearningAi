package domain

type RuleBasedLearningStyleClassifier struct{}

func NewRuleBasedLearningStyleClassifier() *RuleBasedLearningStyleClassifier {
	return &RuleBasedLearningStyleClassifier{}
}

func (c *RuleBasedLearningStyleClassifier) Classify(signals []LearningSignal) map[string]float64 {
	scores := map[string]float64{
		"visual":    0.0,
		"narrative": 0.0,
		"concise":   0.0,
	}

	for _, s := range signals {
		weight := CalculateInteractionScore(s.EngagementScore, s.FeedbackScore, NormalizeTimeSpent(s.TimeSpent))

		switch s.ExplanationType {
		case "visual":
			scores["visual"] += weight
		case "narrative", "analogy":
			scores["narrative"] += weight
		case "concise", "summary":
			scores["concise"] += weight
		}
	}

	return normalizeScores(scores)
}

type RuleBasedInterestClassifier struct{}

func NewRuleBasedInterestClassifier() *RuleBasedInterestClassifier {
	return &RuleBasedInterestClassifier{}
}

func (c *RuleBasedInterestClassifier) Classify(signals []LearningSignal) map[string]float64 {
	scores := make(map[string]float64)

	for _, s := range signals {
		if s.ThemeUsed == "" {
			continue
		}

		weight := CalculateInteractionScore(s.EngagementScore, s.FeedbackScore, NormalizeTimeSpent(s.TimeSpent))
		scores[s.ThemeUsed] += weight
	}

	return normalizeScores(scores)
}

func normalizeScores(scores map[string]float64) map[string]float64 {
	var total float64
	for _, v := range scores {
		total += v
	}

	if total == 0 {
		return scores
	}

	normalized := make(map[string]float64)
	for k, v := range scores {
		normalized[k] = v / total
	}
	return normalized
}
