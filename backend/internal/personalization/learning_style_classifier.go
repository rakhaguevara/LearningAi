package personalization

import (
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LearningStyleClassifier analyzes user behavior signals to determine
// their preferred learning style using a weighted scoring algorithm.
type LearningStyleClassifier struct {
	repo *Repository
	log  *zap.Logger
}

func NewLearningStyleClassifier(repo *Repository, log *zap.Logger) *LearningStyleClassifier {
	return &LearningStyleClassifier{repo: repo, log: log}
}

// signalWeights defines how each signal type contributes to learning style scores.
// Positive values increase the score, negative decrease it.
var signalWeights = map[SignalType]map[LearningStyle]float64{
	SignalIllustrationRequest: {
		StyleVisual:      1.5,
		StyleAuditory:    -0.2,
		StyleReading:     -0.3,
		StyleKinesthetic: 0.3,
	},
	SignalExampleRequest: {
		StyleVisual:      0.5,
		StyleAuditory:    0.3,
		StyleReading:     0.2,
		StyleKinesthetic: 1.2,
	},
	SignalRepetitionRequest: {
		StyleVisual:      0.2,
		StyleAuditory:    1.0,
		StyleReading:     0.5,
		StyleKinesthetic: 0.3,
	},
	SignalFollowUpQuestion: {
		StyleVisual:      0.3,
		StyleAuditory:    0.5,
		StyleReading:     1.0,
		StyleKinesthetic: 0.4,
	},
	SignalAnalogySatisfaction: {
		StyleVisual:      0.8,
		StyleAuditory:    0.6,
		StyleReading:     0.4,
		StyleKinesthetic: 1.0,
	},
	SignalSessionDuration: {
		StyleVisual:      0.3,
		StyleAuditory:    0.3,
		StyleReading:     0.8,
		StyleKinesthetic: 0.2,
	},
	SignalResponseEngagement: {
		StyleVisual:      0.4,
		StyleAuditory:    0.4,
		StyleReading:     0.6,
		StyleKinesthetic: 0.5,
	},
}

// Classify analyzes behavior signals and returns a learning style profile.
func (c *LearningStyleClassifier) Classify(userID uuid.UUID) (*LearningStyleProfile, error) {
	// Get signals from the last 30 days for classification
	since := time.Now().AddDate(0, 0, -30)
	signals, err := c.repo.GetUserSignals(userID, 500)
	if err != nil {
		return nil, err
	}

	// Filter to recent signals
	var recentSignals []BehaviorSignal
	for _, s := range signals {
		if s.CreatedAt.After(since) {
			recentSignals = append(recentSignals, s)
		}
	}

	// Initialize scores
	scores := map[LearningStyle]float64{
		StyleVisual:      0,
		StyleAuditory:    0,
		StyleReading:     0,
		StyleKinesthetic: 0,
	}

	// Calculate weighted scores
	for _, signal := range recentSignals {
		weights, exists := signalWeights[signal.SignalType]
		if !exists {
			continue
		}

		// Apply recency decay (more recent signals matter more)
		daysSince := time.Since(signal.CreatedAt).Hours() / 24
		recencyFactor := math.Exp(-daysSince / 15) // Half-life of ~15 days

		// Apply signal value as intensity multiplier
		intensityFactor := 0.5 + (signal.Value * 0.5) // Range 0.5-1.0

		for style, weight := range weights {
			scores[style] += weight * recencyFactor * intensityFactor
		}
	}

	// Normalize scores to 0-1 range
	maxScore := 0.0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore > 0 {
		for style := range scores {
			scores[style] = scores[style] / maxScore
		}
	}

	// Determine primary style
	primaryStyle := StyleAdaptive
	highestScore := 0.0
	for style, score := range scores {
		if score > highestScore {
			highestScore = score
			primaryStyle = style
		}
	}

	// Calculate confidence based on score differentiation and sample size
	confidence := c.calculateConfidence(scores, len(recentSignals))

	// If confidence is too low, default to adaptive
	if confidence < 0.3 {
		primaryStyle = StyleAdaptive
	}

	profile := &LearningStyleProfile{
		UserID:           userID,
		PrimaryStyle:     primaryStyle,
		VisualScore:      scores[StyleVisual],
		AuditoryScore:    scores[StyleAuditory],
		ReadingScore:     scores[StyleReading],
		KinestheticScore: scores[StyleKinesthetic],
		Confidence:       confidence,
		SampleSize:       len(recentSignals),
		LastCalculatedAt: time.Now(),
	}

	// Persist the profile
	if err := c.repo.UpsertLearningStyleProfile(profile); err != nil {
		c.log.Error("failed to save learning style profile", zap.Error(err))
	}

	c.log.Info("classified learning style",
		zap.String("user_id", userID.String()),
		zap.String("primary_style", string(primaryStyle)),
		zap.Float64("confidence", confidence),
		zap.Int("sample_size", len(recentSignals)),
	)

	return profile, nil
}

// calculateConfidence determines how confident we are in the classification.
func (c *LearningStyleClassifier) calculateConfidence(scores map[LearningStyle]float64, sampleSize int) float64 {
	// Factor 1: Sample size (more data = more confidence)
	// Sigmoid function: approaches 1 as sample size increases
	sampleConfidence := 1 - math.Exp(-float64(sampleSize)/50)

	// Factor 2: Score differentiation (clear winner = more confidence)
	var sorted []float64
	for _, s := range scores {
		sorted = append(sorted, s)
	}
	// Simple sort for 4 elements
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] > sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	differentiationConfidence := 0.0
	if len(sorted) >= 2 && sorted[0] > 0 {
		// How much does the top score dominate the second?
		differentiationConfidence = (sorted[0] - sorted[1]) / sorted[0]
	}

	// Combine factors (weighted average)
	confidence := (sampleConfidence * 0.4) + (differentiationConfidence * 0.6)

	// Clamp to 0-1
	if confidence > 1 {
		confidence = 1
	}
	if confidence < 0 {
		confidence = 0
	}

	return confidence
}

// GetStyleRecommendations returns teaching recommendations based on learning style.
func (c *LearningStyleClassifier) GetStyleRecommendations(style LearningStyle) StyleRecommendations {
	recommendations := map[LearningStyle]StyleRecommendations{
		StyleVisual: {
			PreferredFormats: []string{"diagrams", "charts", "infographics", "videos", "color-coded text"},
			ExplanationStyle: "Use visual metaphors, describe spatial relationships, include imagery",
			AnalogyApproach:  "Draw from visual media: movies, art, architecture, nature scenes",
			ExampleTypes:     []string{"illustrated examples", "flowcharts", "mind maps"},
			AvoidPatterns:    []string{"long unbroken text", "audio-only content"},
		},
		StyleAuditory: {
			PreferredFormats: []string{"explanations", "discussions", "verbal walkthroughs", "mnemonics"},
			ExplanationStyle: "Use rhythm and patterns in explanation, repeat key concepts",
			AnalogyApproach:  "Draw from music, conversations, podcasts, storytelling",
			ExampleTypes:     []string{"verbal scenarios", "dialogue-based examples"},
			AvoidPatterns:    []string{"dense visual diagrams without explanation"},
		},
		StyleReading: {
			PreferredFormats: []string{"detailed text", "bullet points", "documentation", "written examples"},
			ExplanationStyle: "Structured, logical progression with clear headers and sections",
			AnalogyApproach:  "Draw from literature, articles, written narratives",
			ExampleTypes:     []string{"code samples", "written case studies", "step-by-step guides"},
			AvoidPatterns:    []string{"overly simplified content", "too many images"},
		},
		StyleKinesthetic: {
			PreferredFormats: []string{"interactive examples", "hands-on exercises", "real-world applications"},
			ExplanationStyle: "Focus on practical application, use action verbs, physical analogies",
			AnalogyApproach:  "Draw from sports, physical activities, hands-on hobbies, building/crafting",
			ExampleTypes:     []string{"try-it-yourself exercises", "real-world scenarios", "simulations"},
			AvoidPatterns:    []string{"purely theoretical content", "passive reading"},
		},
		StyleAdaptive: {
			PreferredFormats: []string{"mixed media", "varied approaches", "balanced content"},
			ExplanationStyle: "Combine visual, verbal, and practical elements",
			AnalogyApproach:  "Use diverse analogy sources based on topic",
			ExampleTypes:     []string{"varied example types", "multi-modal content"},
			AvoidPatterns:    []string{"single-format content only"},
		},
	}

	if rec, exists := recommendations[style]; exists {
		return rec
	}
	return recommendations[StyleAdaptive]
}

type StyleRecommendations struct {
	PreferredFormats []string `json:"preferred_formats"`
	ExplanationStyle string   `json:"explanation_style"`
	AnalogyApproach  string   `json:"analogy_approach"`
	ExampleTypes     []string `json:"example_types"`
	AvoidPatterns    []string `json:"avoid_patterns"`
}
