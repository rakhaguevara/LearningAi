package user

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

func (h *Handler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Err(c, apperr.NewUnauthorized("user not authenticated"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Err(c, apperr.NewBadRequest("invalid user id"))
		return
	}

	profile, err := h.service.GetProfile(userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, profile)
}
