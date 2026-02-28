package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/adaptive-ai-learn/backend/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GoogleAuth(c *gin.Context) {
	var payload GoogleTokenPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.service.AuthenticateWithGoogle(payload.Token)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}
