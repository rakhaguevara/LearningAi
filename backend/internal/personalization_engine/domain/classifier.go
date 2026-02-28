package domain

// LearningStyleClassifier computes a weighted distribution of preferred learning styles
// based on recent behavioral signals.
type LearningStyleClassifier interface {
	Classify(signals []LearningSignal) map[string]float64
}

// InterestClassifier computes a weighted distribution of preferred content themes
// (anime, sports, tech, etc.) based on recent behavioral signals.
type InterestClassifier interface {
	Classify(signals []LearningSignal) map[string]float64
}
