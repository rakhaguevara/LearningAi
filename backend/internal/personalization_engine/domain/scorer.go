package domain

// CalculateInteractionScore applies a weighted formula to behavioral signals.
// User requested definition:
// FinalScore = (EngagementScore * 0.4) + (FeedbackScore * 0.3) + (TimeSpentWeight * 0.3)
func CalculateInteractionScore(engagementScore, feedbackScore, timeSpentWeight float64) float64 {
	return (engagementScore * 0.4) + (feedbackScore * 0.3) + (timeSpentWeight * 0.3)
}

// NormalizeTimeSpent normalizes raw seconds into a 0.0 - 1.0 weight.
// E.g., we might say 120 seconds is "perfect" (1.0), and anything over is capped.
func NormalizeTimeSpent(seconds int) float64 {
	const maxOptimalSeconds = 120.0
	val := float64(seconds) / maxOptimalSeconds
	if val > 1.0 {
		return 1.0
	}
	return val
}

// ApplyTimeDecay reduces a score based on how many days have passed since it was updated.
// This prevents stale preferences from dominating the profile.
func ApplyTimeDecay(currentScore float64, daysOld int) float64 {
	decayFactor := 0.95 // 5% decay per day
	for i := 0; i < daysOld; i++ {
		currentScore *= decayFactor
	}
	return currentScore
}
