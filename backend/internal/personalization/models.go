package personalization

import (
	"time"

	"github.com/google/uuid"
)

// ─── Learning Styles ────────────────────────────────────────────────────────

type LearningStyle string

const (
	StyleVisual      LearningStyle = "visual"
	StyleAuditory    LearningStyle = "auditory"
	StyleReading     LearningStyle = "reading"
	StyleKinesthetic LearningStyle = "kinesthetic"
	StyleAdaptive    LearningStyle = "adaptive"
)

// ─── Behavior Signal Types ──────────────────────────────────────────────────

type SignalType string

const (
	SignalExplanationRequest  SignalType = "explanation_request"
	SignalIllustrationRequest SignalType = "illustration_request"
	SignalFollowUpQuestion    SignalType = "follow_up_question"
	SignalTopicSwitch         SignalType = "topic_switch"
	SignalDifficultyFeedback  SignalType = "difficulty_feedback"
	SignalSessionDuration     SignalType = "session_duration"
	SignalResponseEngagement  SignalType = "response_engagement"
	SignalInterestIndication  SignalType = "interest_indication"
	SignalRepetitionRequest   SignalType = "repetition_request"
	SignalExampleRequest      SignalType = "example_request"
	SignalAnalogySatisfaction SignalType = "analogy_satisfaction"
)

// ─── Behavior Signal ────────────────────────────────────────────────────────

type BehaviorSignal struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	UserID     uuid.UUID              `json:"user_id" db:"user_id"`
	SessionID  *uuid.UUID             `json:"session_id,omitempty" db:"session_id"`
	SignalType SignalType             `json:"signal_type" db:"signal_type"`
	Value      float64                `json:"value" db:"value"`
	Context    map[string]interface{} `json:"context" db:"context"`
	Topic      string                 `json:"topic,omitempty" db:"topic"`
	Subject    string                 `json:"subject,omitempty" db:"subject"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}

// ─── Learning Style Profile ─────────────────────────────────────────────────

type LearningStyleProfile struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserID           uuid.UUID     `json:"user_id" db:"user_id"`
	PrimaryStyle     LearningStyle `json:"primary_style" db:"primary_style"`
	VisualScore      float64       `json:"visual_score" db:"visual_score"`
	AuditoryScore    float64       `json:"auditory_score" db:"auditory_score"`
	ReadingScore     float64       `json:"reading_score" db:"reading_score"`
	KinestheticScore float64       `json:"kinesthetic_score" db:"kinesthetic_score"`
	Confidence       float64       `json:"confidence" db:"confidence"`
	SampleSize       int           `json:"sample_size" db:"sample_size"`
	LastCalculatedAt time.Time     `json:"last_calculated_at" db:"last_calculated_at"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

// ─── Interest Profile ───────────────────────────────────────────────────────

type InterestProfile struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	UserID         uuid.UUID        `json:"user_id" db:"user_id"`
	Interests      []InterestWeight `json:"interests"`
	TopCategories  []string         `json:"top_categories"`
	AnalogySources []string         `json:"analogy_sources"`
	LastUpdatedAt  time.Time        `json:"last_updated_at" db:"last_updated_at"`
}

type InterestWeight struct {
	Tag        string  `json:"tag"`
	Category   string  `json:"category"`
	Weight     float64 `json:"weight"`
	Engagement float64 `json:"engagement"`
	Recency    float64 `json:"recency"`
}

// ─── Engagement Metrics ─────────────────────────────────────────────────────

type EngagementMetrics struct {
	UserID              uuid.UUID `json:"user_id"`
	TotalSessions       int       `json:"total_sessions"`
	AvgSessionDuration  float64   `json:"avg_session_duration_sec"`
	AvgQuestionsPerSess float64   `json:"avg_questions_per_session"`
	CompletionRate      float64   `json:"completion_rate"`
	ReturnRate          float64   `json:"return_rate"`
	DepthScore          float64   `json:"depth_score"`
	EngagementLevel     string    `json:"engagement_level"`
}

// ─── Personalization Profile (Combined Output) ──────────────────────────────

type PersonalizationProfile struct {
	UserID              uuid.UUID            `json:"user_id"`
	LearningStyle       LearningStyleProfile `json:"learning_style"`
	Interests           InterestProfile      `json:"interests"`
	Engagement          EngagementMetrics    `json:"engagement"`
	PreferredComplexity string               `json:"preferred_complexity"`
	PreferredTone       string               `json:"preferred_tone"`
	AnalogyDomains      []string             `json:"analogy_domains"`
	AdaptivePrompt      string               `json:"adaptive_prompt"`
}

// ─── Tone Configuration ─────────────────────────────────────────────────────

type ToneConfig struct {
	Formality     float64 `json:"formality"`     // 0 = casual, 1 = formal
	Enthusiasm    float64 `json:"enthusiasm"`    // 0 = neutral, 1 = excited
	Technicality  float64 `json:"technicality"`  // 0 = simple, 1 = technical
	Verbosity     float64 `json:"verbosity"`     // 0 = concise, 1 = detailed
	Encouragement float64 `json:"encouragement"` // 0 = neutral, 1 = supportive
}
