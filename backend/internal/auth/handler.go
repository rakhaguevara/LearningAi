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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.service.Register(req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}
