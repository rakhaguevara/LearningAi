package onboarding

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetStatus handles GET /onboarding/status
func (h *Handler) GetStatus(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		response.Err(c, apperr.NewUnauthorized("invalid user context"))
		return
	}

	completed, err := h.service.GetStatus(userID)
	if err != nil {
		response.Err(c, apperr.NewInternal(err.Error()))
		return
	}

	response.OK(c, OnboardingStatusResponse{ProfileCompleted: completed})
}

// Submit handles POST /onboarding/submit
func (h *Handler) Submit(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		response.Err(c, apperr.NewUnauthorized("invalid user context"))
		return
	}

	var req OnboardingSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, apperr.NewBadRequest(err.Error()))
		return
	}

	if err := h.service.Submit(userID, req); err != nil {
		response.Err(c, apperr.NewInternal(err.Error()))
		return
	}

	response.OK(c, gin.H{"message": "onboarding completed successfully"})
}

// UpdateLearning handles PUT /profile/update-learning
func (h *Handler) UpdateLearning(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		response.Err(c, apperr.NewUnauthorized("invalid user context"))
		return
	}

	var req UpdateLearningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, apperr.NewBadRequest(err.Error()))
		return
	}

	if err := h.service.UpdateLearning(userID, req); err != nil {
		response.Err(c, apperr.NewInternal(err.Error()))
		return
	}

	response.OK(c, gin.H{"message": "learning profile updated"})
}

// extractUserID reads the user_id from the Gin context (set by the auth middleware as "userID").
func extractUserID(c *gin.Context) (uuid.UUID, error) {
	raw, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, apperr.NewUnauthorized("missing user context")
	}
	switch v := raw.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, apperr.NewUnauthorized("invalid user context type")
	}
}
