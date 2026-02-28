package personalization

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ScoringEngine combines all classifiers to produce a unified personalization profile.
type ScoringEngine struct {
	repo               *Repository
	styleClassifier    *LearningStyleClassifier
	interestClassifier *InterestClassifier
	log                *zap.Logger
}

func NewScoringEngine(
	repo *Repository,
	styleClassifier *LearningStyleClassifier,
	interestClassifier *InterestClassifier,
	log *zap.Logger,
) *ScoringEngine {
	return &ScoringEngine{
		repo:               repo,
		styleClassifier:    styleClassifier,
		interestClassifier: interestClassifier,
		log:                log,
	}
}

// BuildProfile constructs a complete personalization profile for a user.
func (e *ScoringEngine) BuildProfile(userID uuid.UUID) (*PersonalizationProfile, error) {
	// Get learning style classification
	styleProfile, err := e.styleClassifier.Classify(userID)
	if err != nil {
		e.log.Warn("learning style classification failed, using defaults", zap.Error(err))
		styleProfile = &LearningStyleProfile{
			UserID:       userID,
			PrimaryStyle: StyleAdaptive,
			Confidence:   0,
		}
	}

	// Get interest classification
	interestProfile, err := e.interestClassifier.Classify(userID)
	if err != nil {
		e.log.Warn("interest classification failed, using defaults", zap.Error(err))
		interestProfile = &InterestProfile{
			UserID:         userID,
			Interests:      []InterestWeight{},
			TopCategories:  []string{},
			AnalogySources: []string{},
		}
	}

	// Get engagement metrics
	engagement, err := e.repo.GetEngagementMetrics(userID)
	if err != nil {
		e.log.Warn("engagement metrics failed, using defaults", zap.Error(err))
		engagement = &EngagementMetrics{
			UserID:          userID,
			EngagementLevel: "low",
		}
	}

	// Determine preferred complexity based on engagement
	preferredComplexity := e.determineComplexity(engagement, styleProfile)

	// Determine preferred tone
	preferredTone := e.determineTone(engagement, styleProfile)

	// Build analogy domains from interests
	analogyDomains := interestProfile.AnalogySources
	if len(analogyDomains) == 0 {
		analogyDomains = []string{"general knowledge", "everyday examples"}
	}

	// Generate adaptive system prompt
	adaptivePrompt := e.generateAdaptivePrompt(styleProfile, interestProfile, preferredComplexity, preferredTone)

	profile := &PersonalizationProfile{
		UserID:              userID,
		LearningStyle:       *styleProfile,
		Interests:           *interestProfile,
		Engagement:          *engagement,
		PreferredComplexity: preferredComplexity,
		PreferredTone:       preferredTone,
		AnalogyDomains:      analogyDomains,
		AdaptivePrompt:      adaptivePrompt,
	}

	e.log.Info("built personalization profile",
		zap.String("user_id", userID.String()),
		zap.String("style", string(styleProfile.PrimaryStyle)),
		zap.String("complexity", preferredComplexity),
		zap.String("tone", preferredTone),
		zap.Int("interests_count", len(interestProfile.Interests)),
	)

	return profile, nil
}

// determineComplexity calculates preferred explanation complexity.
func (e *ScoringEngine) determineComplexity(engagement *EngagementMetrics, style *LearningStyleProfile) string {
	score := 0.0

	switch engagement.EngagementLevel {
	case "high":
		score += 0.6
	case "medium":
		score += 0.4
	default:
		score += 0.2
	}

	if engagement.AvgSessionDuration > 900 {
		score += 0.2
	} else if engagement.AvgSessionDuration > 300 {
		score += 0.1
	}

	if style.PrimaryStyle == StyleReading {
		score += 0.15
	}

	if style.Confidence > 0.7 {
		score += 0.1
	}

	if score >= 0.7 {
		return "advanced"
	} else if score >= 0.4 {
		return "intermediate"
	}
	return "beginner"
}

// determineTone calculates preferred explanation tone.
func (e *ScoringEngine) determineTone(engagement *EngagementMetrics, style *LearningStyleProfile) string {
	if style.PrimaryStyle == StyleKinesthetic {
		return "casual_engaging"
	}

	if style.PrimaryStyle == StyleReading && style.Confidence > 0.5 {
		return "structured_formal"
	}

	if engagement.EngagementLevel == "high" {
		return "professional_encouraging"
	}

	return "friendly_accessible"
}

// generateAdaptivePrompt creates a system prompt for the AI based on profile.
func (e *ScoringEngine) generateAdaptivePrompt(
	style *LearningStyleProfile,
	interests *InterestProfile,
	complexity string,
	tone string,
) string {
	var sb strings.Builder

	sb.WriteString("You are an adaptive AI tutor. ")

	styleRecs := e.styleClassifier.GetStyleRecommendations(style.PrimaryStyle)
	sb.WriteString(fmt.Sprintf("\n\nLEARNING STYLE: The user is primarily a %s learner. ", style.PrimaryStyle))
	sb.WriteString(styleRecs.ExplanationStyle)
	sb.WriteString(fmt.Sprintf(" Prefer these formats: %s. ", strings.Join(styleRecs.PreferredFormats, ", ")))
	sb.WriteString(fmt.Sprintf("Avoid: %s.", strings.Join(styleRecs.AvoidPatterns, ", ")))

	if len(interests.Interests) > 0 {
		var topInterests []string
		for i := 0; i < len(interests.Interests) && i < 5; i++ {
			topInterests = append(topInterests, interests.Interests[i].Tag)
		}
		sb.WriteString(fmt.Sprintf("\n\nINTERESTS: The user is interested in: %s. ", strings.Join(topInterests, ", ")))
		sb.WriteString("Draw analogies, metaphors, and examples from these domains whenever possible. ")
		sb.WriteString("Make explanations relatable by connecting concepts to these interests.")
	}

	if len(interests.AnalogySources) > 0 {
		maxSources := len(interests.AnalogySources)
		if maxSources > 8 {
			maxSources = 8
		}
		sb.WriteString(fmt.Sprintf("\n\nANALOGY DOMAINS: Use references from: %s.", strings.Join(interests.AnalogySources[:maxSources], ", ")))
	}

	sb.WriteString("\n\nCOMPLEXITY LEVEL: ")
	switch complexity {
	case "advanced":
		sb.WriteString("The user is comfortable with advanced content. Use technical terminology, go into depth, and don't oversimplify.")
	case "intermediate":
		sb.WriteString("Balance accessibility with depth. Introduce technical terms with brief explanations.")
	default:
		sb.WriteString("Keep explanations simple and accessible. Use everyday language. Break down complex ideas.")
	}

	sb.WriteString("\n\nTONE: ")
	switch tone {
	case "casual_engaging":
		sb.WriteString("Be casual, energetic, and action-oriented. Use active voice.")
	case "structured_formal":
		sb.WriteString("Be organized and precise. Use clear structure with logical progression.")
	case "professional_encouraging":
		sb.WriteString("Be professional but warm. Acknowledge the user's capability while offering support.")
	default:
		sb.WriteString("Be friendly, warm, and accessible. Use conversational language.")
	}

	sb.WriteString("\n\nREMEMBER: Every explanation should feel personally crafted for this learner.")

	return sb.String()
}

// CalculatePersonalizationScore returns a 0-100 score indicating personalization depth.
func (e *ScoringEngine) CalculatePersonalizationScore(profile *PersonalizationProfile) int {
	score := 0.0

	score += profile.LearningStyle.Confidence * 25

	sampleScore := float64(profile.LearningStyle.SampleSize) / 100 * 15
	if sampleScore > 15 {
		sampleScore = 15
	}
	score += sampleScore

	interestScore := float64(len(profile.Interests.Interests)) / 10 * 20
	if interestScore > 20 {
		interestScore = 20
	}
	score += interestScore

	categoryScore := float64(len(profile.Interests.TopCategories)) / 5 * 15
	if categoryScore > 15 {
		categoryScore = 15
	}
	score += categoryScore

	switch profile.Engagement.EngagementLevel {
	case "high":
		score += 25
	case "medium":
		score += 15
	default:
		score += 5
	}

	return int(score)
}
