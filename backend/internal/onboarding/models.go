package onboarding

import (
	"time"

	"github.com/google/uuid"
)

// OnboardingSubmitRequest contains answers to all 10 onboarding questions.
type OnboardingSubmitRequest struct {
	LearningStyle       string   `json:"learning_style" binding:"required,oneof=visual step_by_step stories hands_on"`
	LongContentBehavior string   `json:"long_content_behavior" binding:"required,oneof=summarize keep_going break_up give_examples"`
	ExplanationFormat   string   `json:"explanation_format" binding:"required,oneof=bullet_points paragraphs code_examples diagrams"`
	InterestThemes      []string `json:"interest_themes" binding:"required,min=1,dive,oneof=tech science arts sports gaming business music"`
	AnalogyTheme        string   `json:"analogy_theme" binding:"required,oneof=sports gaming cooking movies nature"`
	DepthPreference     string   `json:"depth_preference" binding:"required,oneof=beginner_overview deep_dive concept_plus_examples expert"`
	AIRetryPreference   string   `json:"ai_retry_preference" binding:"required,oneof=rephrase try_simpler ask_examples give_up"`
	StudyFocus          string   `json:"study_focus" binding:"required,max=200"`
	FileUploadHabit     string   `json:"file_upload_habit" binding:"required,oneof=yes_always sometimes rarely no"`
	LearningGoal        string   `json:"learning_goal" binding:"required,oneof=pass_exams learn_for_fun career_change build_projects"`
}

// UpdateLearningRequest allows partial post-onboarding updates to learning preferences.
type UpdateLearningRequest struct {
	LearningStyle     *string  `json:"learning_style"`
	DepthPreference   *string  `json:"depth_preference"`
	InterestThemes    []string `json:"interest_themes"`
	ExplanationFormat *string  `json:"explanation_format"`
	LearningGoal      *string  `json:"learning_goal"`
}

// OnboardingStatusResponse reports whether the user has completed onboarding.
type OnboardingStatusResponse struct {
	ProfileCompleted bool `json:"profile_completed"`
}

// UserLearningProfile maps to the user_learning_profiles table.
type UserLearningProfile struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Preferences []byte    `json:"preferences"` // raw JSONB
	Meta        []byte    `json:"meta"`        // raw JSONB
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PreferencesPayload is the canonically serialized form stored in user_learning_profiles.preferences.
type PreferencesPayload struct {
	LearningStyle       string   `json:"learning_style"`
	LongContentBehavior string   `json:"long_content_behavior"`
	ExplanationFormat   string   `json:"explanation_format"`
	InterestThemes      []string `json:"interest_themes"`
	AnalogyTheme        string   `json:"analogy_theme"`
	DepthPreference     string   `json:"depth_preference"`
	AIRetryPreference   string   `json:"ai_retry_preference"`
	StudyFocus          string   `json:"study_focus"`
	FileUploadHabit     string   `json:"file_upload_habit"`
	LearningGoal        string   `json:"learning_goal"`
}
