package personalization

import (
	"fmt"
	"strings"
)

// ToneGenerator produces tone configuration and instructions based on user profile.
type ToneGenerator struct{}

func NewToneGenerator() *ToneGenerator {
	return &ToneGenerator{}
}

// GenerateToneConfig produces numeric tone parameters from a profile.
func (g *ToneGenerator) GenerateToneConfig(profile *PersonalizationProfile) ToneConfig {
	config := ToneConfig{
		Formality:     0.5,
		Enthusiasm:    0.5,
		Technicality:  0.5,
		Verbosity:     0.5,
		Encouragement: 0.5,
	}

	// Adjust based on preferred tone
	switch profile.PreferredTone {
	case "casual_engaging":
		config.Formality = 0.2
		config.Enthusiasm = 0.8
		config.Encouragement = 0.7
	case "structured_formal":
		config.Formality = 0.8
		config.Enthusiasm = 0.4
		config.Verbosity = 0.7
	case "professional_encouraging":
		config.Formality = 0.6
		config.Enthusiasm = 0.6
		config.Encouragement = 0.8
	case "friendly_accessible":
		config.Formality = 0.3
		config.Enthusiasm = 0.6
		config.Encouragement = 0.7
	}

	// Adjust based on complexity
	switch profile.PreferredComplexity {
	case "advanced":
		config.Technicality = 0.8
		config.Verbosity = 0.7
	case "intermediate":
		config.Technicality = 0.5
		config.Verbosity = 0.5
	case "beginner":
		config.Technicality = 0.2
		config.Verbosity = 0.6
		config.Encouragement += 0.1
	}

	// Adjust based on learning style
	switch profile.LearningStyle.PrimaryStyle {
	case StyleVisual:
		config.Verbosity -= 0.1 // Prefer concise with visuals
	case StyleReading:
		config.Verbosity += 0.15 // Prefer detailed text
		config.Formality += 0.1
	case StyleKinesthetic:
		config.Enthusiasm += 0.1
		config.Formality -= 0.1
	case StyleAuditory:
		config.Verbosity += 0.1 // Prefer explanatory
	}

	// Adjust based on engagement
	switch profile.Engagement.EngagementLevel {
	case "high":
		config.Technicality += 0.1
	case "low":
		config.Encouragement += 0.15
		config.Enthusiasm += 0.1
	}

	// Clamp all values to 0-1
	config.Formality = clamp(config.Formality)
	config.Enthusiasm = clamp(config.Enthusiasm)
	config.Technicality = clamp(config.Technicality)
	config.Verbosity = clamp(config.Verbosity)
	config.Encouragement = clamp(config.Encouragement)

	return config
}

// GenerateToneInstructions produces natural language instructions for the AI.
func (g *ToneGenerator) GenerateToneInstructions(config ToneConfig) string {
	var instructions []string

	// Formality
	if config.Formality < 0.3 {
		instructions = append(instructions, "Use casual, conversational language. Contractions are welcome. Feel free to use colloquialisms.")
	} else if config.Formality > 0.7 {
		instructions = append(instructions, "Maintain a professional, formal tone. Use proper grammar and complete sentences.")
	} else {
		instructions = append(instructions, "Balance professionalism with approachability.")
	}

	// Enthusiasm
	if config.Enthusiasm > 0.7 {
		instructions = append(instructions, "Be energetic and enthusiastic! Show excitement about the topic.")
	} else if config.Enthusiasm < 0.3 {
		instructions = append(instructions, "Keep a calm, measured tone. Avoid excessive exclamation.")
	}

	// Technicality
	if config.Technicality > 0.7 {
		instructions = append(instructions, "Use technical terminology confidently. The user can handle domain-specific language.")
	} else if config.Technicality < 0.3 {
		instructions = append(instructions, "Avoid jargon. If technical terms are necessary, explain them in simple words.")
	} else {
		instructions = append(instructions, "Introduce technical terms but provide brief explanations when first used.")
	}

	// Verbosity
	if config.Verbosity > 0.7 {
		instructions = append(instructions, "Provide detailed, thorough explanations. Include context and background.")
	} else if config.Verbosity < 0.3 {
		instructions = append(instructions, "Be concise and to the point. Focus on key information.")
	}

	// Encouragement
	if config.Encouragement > 0.7 {
		instructions = append(instructions, "Be supportive and encouraging. Acknowledge progress and effort. Use positive reinforcement.")
	}

	return strings.Join(instructions, " ")
}

// GenerateStylePromptFragment produces a prompt fragment for a specific learning style.
func (g *ToneGenerator) GenerateStylePromptFragment(style LearningStyle, interests []string) string {
	var fragments []string

	switch style {
	case StyleVisual:
		fragments = append(fragments,
			"Describe concepts visually - use spatial relationships, colors, shapes, and imagery.",
			"When possible, suggest diagrams or visual representations.",
			"Use phrases like 'picture this', 'imagine seeing', 'visualize'.",
		)
	case StyleAuditory:
		fragments = append(fragments,
			"Explain concepts as if speaking aloud - use rhythm, repetition, and verbal patterns.",
			"Include mnemonics or memorable phrases when helpful.",
			"Structure explanations like a conversation or lecture.",
		)
	case StyleReading:
		fragments = append(fragments,
			"Provide well-structured, written explanations with clear organization.",
			"Use bullet points, numbered lists, and hierarchical structure.",
			"Include references and suggest further reading when relevant.",
		)
	case StyleKinesthetic:
		fragments = append(fragments,
			"Focus on practical, hands-on applications.",
			"Use action verbs and physical metaphors.",
			"Include 'try this' exercises or thought experiments.",
			"Relate concepts to real-world actions and experiences.",
		)
	default:
		fragments = append(fragments,
			"Use a balanced mix of visual, verbal, and practical explanations.",
			"Adapt your approach based on the complexity of the topic.",
		)
	}

	// Add interest-based instructions
	if len(interests) > 0 {
		interestList := strings.Join(interests[:minInt(3, len(interests))], ", ")
		fragments = append(fragments,
			fmt.Sprintf("Draw examples and analogies from: %s.", interestList),
		)
	}

	return strings.Join(fragments, " ")
}

// AdaptResponseDynamically adjusts a response based on real-time signals.
func (g *ToneGenerator) AdaptResponseDynamically(
	originalResponse string,
	feedbackSignal SignalType,
	feedbackValue float64,
) string {
	// This would be called after receiving user feedback to suggest adjustments
	// In a real implementation, this could trigger a re-generation with adjusted parameters

	adjustments := map[SignalType]string{
		SignalRepetitionRequest: "The user requested repetition - they may need simpler or slower explanation.",
		SignalDifficultyFeedback: func() string {
			if feedbackValue < 0.5 {
				return "User found content too difficult - simplify and add more foundational context."
			}
			return "User found content too easy - can increase complexity and depth."
		}(),
		SignalExampleRequest:      "User wants more examples - include additional practical applications.",
		SignalIllustrationRequest: "User wants visual aids - describe concepts more visually or suggest diagrams.",
	}

	if adjustment, exists := adjustments[feedbackSignal]; exists {
		return fmt.Sprintf("[Adaptation Note: %s]\n\n%s", adjustment, originalResponse)
	}

	return originalResponse
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
