package ai

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/adaptive-ai-learn/backend/internal/common/response"
)

// Handler exposes all AI workspace endpoints.
type Handler struct {
	svc *AIService
	log *zap.Logger
}

func NewHandler(svc *AIService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// extractUserID reads the user UUID from Gin context (set by Auth middleware as "userID").
func extractUserID(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("userID")
	if !exists {
		response.ErrStatus(c, http.StatusUnauthorized, fmt.Errorf("user_id not in context"))
		return uuid.Nil, false
	}
	uid, err := uuid.Parse(fmt.Sprintf("%v", raw))
	if err != nil {
		response.ErrStatus(c, http.StatusUnauthorized, fmt.Errorf("invalid user_id in context"))
		return uuid.Nil, false
	}
	return uid, true
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /ai/ask
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) Ask(c *gin.Context) {
	var req AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		return
	}

	// Sanitise question
	req.Question = strings.TrimSpace(req.Question)
	if len(req.Question) == 0 {
		response.ErrStatus(c, http.StatusBadRequest, fmt.Errorf("question cannot be empty"))
		return
	}

	result, err := h.svc.Ask(c.Request.Context(), userID, req)
	if err != nil {
		h.log.Error("Ask failed", zap.String("user", userID.String()), zap.Error(err))
		response.Err(c, err)
		return
	}

	if result.NeedsFormat {
		response.OK(c, result)
		return
	}

	response.OK(c, result)
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /ai/upload
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) UploadFile(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.ErrStatus(c, http.StatusBadRequest, fmt.Errorf("no file provided: %w", err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.ErrStatus(c, http.StatusBadRequest, fmt.Errorf("opening file: %w", err))
		return
	}
	defer file.Close()

	doc, err := h.svc.fileParser.ParseUpload(c.Request.Context(), fileHeader, file)
	if err != nil {
		response.ErrStatus(c, http.StatusUnprocessableEntity, err)
		return
	}

	chunks, err := h.svc.ragEngine.StoreChunks(c.Request.Context(), userID, doc.Text, doc.Source)
	if err != nil {
		h.log.Error("storing chunks failed", zap.Error(err))
		response.Err(c, err)
		return
	}

	response.OK(c, gin.H{
		"source":     doc.Source,
		"file_type":  doc.FileType,
		"word_count": doc.WordCount,
		"chunks":     chunks,
		"message":    fmt.Sprintf("Successfully processed and indexed '%s' (%d chunks)", doc.Source, chunks),
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// GET /ai/sources
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) GetSources(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		return
	}

	sources, err := h.svc.ragEngine.ListUserSources(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, gin.H{"sources": sources})
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /ai/generate-ppt
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) GeneratePPT(c *gin.Context) {
	var req struct {
		Topic   string `json:"topic" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		return
	}

	result, err := h.svc.GeneratePPT(c.Request.Context(), userID, req.Topic, req.Content, false)
	if err != nil {
		h.log.Error("PPT generation failed", zap.Error(err))
		response.Err(c, err)
		return
	}

	response.OK(c, gin.H{
		"file_name":   result.FileName,
		"slide_count": result.SlideCount,
		"download_url": fmt.Sprintf("/ai/download/ppt/%s/%s",
			userID.String(), result.FileName),
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// GET /ai/download/ppt/:user_id/:filename
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) DownloadPPT(c *gin.Context) {
	userIDParam := c.Param("user_id")
	filename := c.Param("filename")

	// Security: sanitise path segments
	if strings.Contains(filename, "..") || strings.Contains(userIDParam, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	// Validate user_id matches authenticated user
	userID, ok := extractUserID(c)
	if !ok {
		return
	}
	if userID.String() != userIDParam {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	safeName := sanitisePathSegment(filename)
	path := filepath.Join("/tmp/ailearn/ppt", sanitisePathSegment(userIDParam), safeName)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.File(path)
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /ai/generate-audio
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) GenerateAudio(c *gin.Context) {
	var req struct {
		Text  string `json:"text" binding:"required,min=1"`
		Voice string `json:"voice,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		return
	}

	result, err := h.svc.GenerateAudio(c.Request.Context(), userID, req.Text, req.Voice)
	if err != nil {
		h.log.Error("audio generation failed", zap.Error(err))
		response.Err(c, err)
		return
	}

	response.OK(c, gin.H{
		"file_name":    result.FileName,
		"format":       result.Format,
		"duration_sec": result.Duration,
		"download_url": fmt.Sprintf("/ai/download/audio/%s/%s",
			userID.String(), result.FileName),
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// GET /ai/download/audio/:user_id/:filename
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) DownloadAudio(c *gin.Context) {
	userIDParam := c.Param("user_id")
	filename := c.Param("filename")

	if strings.Contains(filename, "..") || strings.Contains(userIDParam, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		return
	}
	if userID.String() != userIDParam {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	path := filepath.Join("/tmp/ailearn/audio",
		sanitisePathSegment(userIDParam),
		sanitisePathSegment(filename))

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := "audio/mpeg"
	if ext == ".wav" {
		mimeType = "audio/wav"
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", mimeType)
	c.File(path)
}

// ──────────────────────────────────────────────────────────────────────────────
// POST /ai/translate
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) Translate(c *gin.Context) {
	var req TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	_, ok := extractUserID(c)
	if !ok {
		return
	}

	result, err := h.svc.Translate(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, result)
}

// ──────────────────────────────────────────────────────────────────────────────
// Legacy endpoints (backward compat)
// ──────────────────────────────────────────────────────────────────────────────

func (h *Handler) Explain(c *gin.Context) {
	var req ExplainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	uid, ok := extractUserID(c)
	if !ok {
		return
	}
	req.UserID = uid.String()

	result, err := h.svc.Explain(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) GenerateIllustration(c *gin.Context) {
	var req IllustrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.svc.GenerateIllustration(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, result)
}

// GetOutputFormats returns the list of available output formats (public utility).
func (h *Handler) GetOutputFormats(c *gin.Context) {
	formats := make([]gin.H, 0, len(OutputFormatLabel))
	for k, v := range OutputFormatLabel {
		formats = append(formats, gin.H{"id": string(k), "label": v})
	}
	response.OK(c, gin.H{"formats": formats})
}

// dummy for strconv usage
var _ = strconv.Itoa
