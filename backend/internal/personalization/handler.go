package personalization

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/common/response"
	"github.com/adaptive-ai-learn/backend/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetProfile returns the full personalization profile for the authenticated user.
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	profile, err := h.service.GetPersonalizationProfile(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, apperr.NewInternal("failed to build personalization profile"))
		return
	}

	score := h.service.scoringEngine.CalculatePersonalizationScore(profile)

	response.OK(c, gin.H{
		"profile":               profile,
		"personalization_score": score,
	})
}

// GetAdaptivePrompt returns the AI system prompt for the user.
func (h *Handler) GetAdaptivePrompt(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	prompt, err := h.service.GetAdaptivePrompt(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, apperr.NewInternal("failed to generate adaptive prompt"))
		return
	}

	response.OK(c, gin.H{"prompt": prompt})
}

// GetLearningStyle returns the user's learning style classification.
func (h *Handler) GetLearningStyle(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	styleProfile, err := h.service.GetLearningStyleProfile(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, apperr.NewInternal("failed to classify learning style"))
		return
	}

	response.OK(c, styleProfile)
}

// GetInterests returns the user's interest profile.
func (h *Handler) GetInterests(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	interestProfile, err := h.service.GetInterestProfile(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, apperr.NewInternal("failed to classify interests"))
		return
	}

	response.OK(c, interestProfile)
}

// RecordSignal records a behavior signal from the client.
func (h *Handler) RecordSignal(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	var req struct {
		SignalType string                 `json:"signal_type" binding:"required"`
		Value      float64                `json:"value"`
		SessionID  *string                `json:"session_id"`
		Topic      string                 `json:"topic"`
		Subject    string                 `json:"subject"`
		Context    map[string]interface{} `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, apperr.NewBadRequest("invalid request body"))
		return
	}

	signal := &BehaviorSignal{
		UserID:     userID,
		SignalType: SignalType(req.SignalType),
		Value:      req.Value,
		Topic:      req.Topic,
		Subject:    req.Subject,
		Context:    req.Context,
	}

	if req.SessionID != nil {
		if sid, err := uuid.Parse(*req.SessionID); err == nil {
			signal.SessionID = &sid
		}
	}

	if err := h.service.RecordBehaviorSignal(c.Request.Context(), signal); err != nil {
		response.Err(c, apperr.NewInternal("failed to record signal"))
		return
	}

	response.OK(c, gin.H{"recorded": true})
}

// AddInterest adds an interest for the user.
func (h *Handler) AddInterest(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	var req struct {
		Tag      string  `json:"tag" binding:"required"`
		Category string  `json:"category"`
		Weight   float64 `json:"weight"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, apperr.NewBadRequest("invalid request body"))
		return
	}

	if req.Weight == 0 {
		req.Weight = 1.0
	}
	if req.Category == "" {
		req.Category = "general"
	}

	if err := h.service.AddUserInterest(c.Request.Context(), userID, req.Tag, req.Category, req.Weight); err != nil {
		response.Err(c, apperr.NewInternal("failed to add interest"))
		return
	}

	response.OK(c, gin.H{"added": true})
}

// RecordFeedback records difficulty feedback.
func (h *Handler) RecordFeedback(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		response.Err(c, err)
		return
	}

	var req struct {
		SessionID    *string `json:"session_id"`
		Feedback     string  `json:"feedback" binding:"required"`
		FeedbackType string  `json:"feedback_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, apperr.NewBadRequest("invalid request body"))
		return
	}

	var sessionID *uuid.UUID
	if req.SessionID != nil {
		if sid, err := uuid.Parse(*req.SessionID); err == nil {
			sessionID = &sid
		}
	}

	if err := h.service.RecordDifficultyFeedback(c.Request.Context(), userID, sessionID, req.Feedback); err != nil {
		response.Err(c, apperr.NewInternal("failed to record feedback"))
		return
	}

	response.OK(c, gin.H{"recorded": true})
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get(middleware.ContextUserID)
	if !exists {
		return uuid.Nil, apperr.NewUnauthorized("user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, apperr.NewBadRequest("invalid user id")
	}

	return userID, nil
}
