package learning

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

func (h *Handler) StartSession(c *gin.Context) {
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

	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.service.StartSession(userID, req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, result)
}
