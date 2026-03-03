package ai

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ──────────────────────────────────────────────────────────────────────────────
// Domain request/response types for the unified AI workspace
// ──────────────────────────────────────────────────────────────────────────────

// AskRequest is the payload for the /ai/ask endpoint.
type AskRequest struct {
	Question     string       `json:"question" binding:"required,min=1,max=4000"`
	OutputFormat OutputFormat `json:"output_format,omitempty"`
	SessionID    string       `json:"session_id,omitempty"`
	TargetLang   string       `json:"target_language,omitempty"` // for translation mode
}

// AskResponse is returned by /ai/ask.
type AskResponse struct {
	Answer           string       `json:"answer"`
	OutputFormat     OutputFormat `json:"output_format"`
	TokensUsed       int          `json:"tokens_used"`
	LatencyMs        int          `json:"latency_ms"`
	NeedsFormat      bool         `json:"needs_format,omitempty"`
	FormatPrompt     string       `json:"format_prompt,omitempty"`
	ContextUsed      bool         `json:"context_used"`
	ContextFound     int          `json:"context_found"`
	DownloadURL      string       `json:"download_url,omitempty"`
	IllustrationURL  string       `json:"illustration_url,omitempty"`
	IsStructuredJSON bool         `json:"is_structured_json"`

	// Typed sub-structs — one is non-nil depending on output_format
	Summary  *SummaryResponse  `json:"summary,omitempty"`
	Detailed *DetailedResponse `json:"detailed,omitempty"`
	Anime    *AnimeResponse    `json:"anime,omitempty"`
	Sports   *SportsResponse   `json:"sports,omitempty"`
	Academic *AcademicResponse `json:"academic,omitempty"`
}

// TranslateRequest carries text + target language.
type TranslateRequest struct {
	Text       string `json:"text" binding:"required,min=1"`
	TargetLang string `json:"target_lang" binding:"required"`
}

// TranslateResponse wraps the translated text.
type TranslateResponse struct {
	Translated string `json:"translated"`
	SourceLang string `json:"source_lang"`
	TokensUsed int    `json:"tokens_used"`
}

// ──────────────────────────────────────────────────────────────────────────────
// AIService — orchestrates all AI workspace features
// ──────────────────────────────────────────────────────────────────────────────

// AIService is the unified service that powers the LearnNow workspace.
type AIService struct {
	client       *QwenClient
	orchestrator *AIOrchestrator
	ragEngine    *RAGEngine
	fileParser   *FileParser
	pptGen       *PPTGenerator
	tts          *TTSService
	imageGen     *ImageGenerator
	topicExt     *TopicExtractor
	db           *sql.DB
	log          *zap.Logger
}

func NewAIService(
	client *QwenClient,
	ragEngine *RAGEngine,
	fileParser *FileParser,
	pptGen *PPTGenerator,
	tts *TTSService,
	imageGen *ImageGenerator,
	db *sql.DB,
	log *zap.Logger,
) *AIService {
	orchestrator := NewAIOrchestrator(client, imageGen, log)
	topicExt := NewTopicExtractor(client, log)
	return &AIService{
		client:       client,
		orchestrator: orchestrator,
		ragEngine:    ragEngine,
		fileParser:   fileParser,
		pptGen:       pptGen,
		tts:          tts,
		imageGen:     imageGen,
		topicExt:     topicExt,
		db:           db,
		log:          log,
	}
}

// Ask is the main AI chat endpoint. It:
//  1. Validates and optionally prompts for output format
//  2. Retrieves the user learning profile
//  3. Retrieves RAG context from uploaded docs
//  4. Builds a dynamic system prompt
//  5. Calls the Qwen API
//  6. Logs the interaction
func (s *AIService) Ask(ctx context.Context, userID uuid.UUID, req AskRequest) (*AskResponse, error) {
	start := time.Now()

	// --- Step 1: Output format gate ---
	if req.OutputFormat == "" {
		return &AskResponse{
			NeedsFormat:  true,
			FormatPrompt: BuildOutputFormatPromptRequest(),
		}, nil
	}

	// --- Step 2: Load user learning profile ---
	profile, err := s.loadUserProfile(ctx, userID)
	if err != nil {
		s.log.Warn("failed to load user profile, using defaults", zap.Error(err))
		profile = defaultProfile()
	}

	// --- Step 3: Retrieve RAG context ---
	ragCtx, chunksFound, err := s.ragEngine.RetrieveContext(ctx, userID, req.Question)
	if err != nil {
		s.log.Warn("RAG retrieval failed, proceeding without context", zap.Error(err))
		ragCtx = ""
	}

	if chunksFound == 0 {
		s.log.Warn("RAG returned zero chunks for query",
			zap.String("user_id", userID.String()),
			zap.String("question", req.Question),
		)
	} else {
		s.log.Info("RAG context injected",
			zap.Int("number_of_chunks_found", chunksFound),
			zap.Int("number_of_chunks_injected", chunksFound),
		)
	}

	// --- Depth Control ---
	lowerQ := strings.ToLower(req.Question)
	if strings.Contains(lowerQ, "aku belum kebayang") || strings.Contains(lowerQ, "belum paham") {
		req.OutputFormat = OutputFormatDetailed
		s.log.Info("depth control triggered: switched to detailed explanation", zap.String("trigger", lowerQ))
	}

	// --- Step 4: Topic Extraction ---
	topic, domain := s.topicExt.Extract(ctx, req.Question)
	s.log.Info("topic extracted", zap.String("topic", topic), zap.String("domain", domain))

	// --- Step 5: Build system prompt ---
	sysPrompt := BuildSystemPrompt(PromptBuilderConfig{
		LearningStyle:    profile.learningStyle,
		DominantInterest: profile.dominantInterest,
		ExplanationDepth: profile.explanationDepth,
		OutputFormat:     req.OutputFormat,
		Topic:            topic,
		Domain:           domain,
		RetrievedContext: ragCtx,
		TargetLanguage:   req.TargetLang,
	})

	// --- Step 5: Route based on output format ---
	var finalAnswer string
	var downloadURL string
	var illustrationURL string
	var tokens int
	isStructuredJSON := false

	switch req.OutputFormat {
	case OutputFormatTranslation:
		// Override prompt exclusively for translation
		sysPrompt = fmt.Sprintf(
			"Translate the following text into %s. Preserve structure and formatting.",
			req.TargetLang,
		)
		chatResp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
			Messages: []ChatMessage{
				{Role: "system", Content: sysPrompt},
				{Role: "user", Content: req.Question},
			},
			MaxTokens: 2000,
		})
		if err != nil {
			return nil, fmt.Errorf("qwen chat translation failed: %w", err)
		}
		finalAnswer = chatResp.Content
		tokens = chatResp.TokensUsed

	case OutputFormatSlides:
		// PPT Generation Route
		res, err := s.GeneratePPT(ctx, userID, req.Question, ragCtx, false)
		if err != nil {
			s.log.Warn("failed to generate JSON slides, retrying once with correction", zap.Error(err))
			res, err = s.GeneratePPT(ctx, userID, req.Question, ragCtx, true)
			if err != nil {
				return nil, fmt.Errorf("failed to generate PPT JSON after 2 attempts: %w", err)
			}
		}
		finalAnswer = fmt.Sprintf("I have created a presentation about %s. You can download the slides below.", req.Question)
		downloadURL = fmt.Sprintf("/ai/download/ppt/%s/%s", userID.String(), res.FileName)
		tokens = res.SlideCount * 150

	case OutputFormatAudio:
		// Audio Script Generation Route
		chatResp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
			Messages: []ChatMessage{
				{Role: "system", Content: sysPrompt},
				{Role: "user", Content: req.Question},
			},
			MaxTokens: 2000,
		})
		if err != nil {
			return nil, fmt.Errorf("qwen chat script failed: %w", err)
		}
		tokens = chatResp.TokensUsed
		finalAnswer = chatResp.Content

		// Pipe straight into TTS
		audioRes, err := s.GenerateAudio(ctx, userID, finalAnswer, "")
		if err != nil {
			s.log.Warn("TTS generation failed but script succeeded", zap.Error(err))
			finalAnswer += "\n\n*(Audio generation failed. Showing script only.)*"
		} else {
			downloadURL = fmt.Sprintf("/ai/download/audio/%s/%s", userID.String(), audioRes.FileName)
		}

	case OutputFormatAnime, OutputFormatSports, OutputFormatSummary, OutputFormatDetailed, OutputFormatAcademic:
		// ── Delegate entirely to AIOrchestrator ─────────────────────────────────
		// The orchestrator handles: text gen → JSON parse/retry/validation → image gen
		orchestrated, orchErr := s.orchestrator.GenerateStructured(ctx, req.OutputFormat, sysPrompt, req.Question, topic)
		if orchErr != nil {
			return nil, fmt.Errorf("orchestrator failed: %w", orchErr)
		}
		// Answer = human-readable main content (not raw JSON)
		if orchestrated.Title != "" {
			finalAnswer = orchestrated.Title
			if orchestrated.MainContent != "" {
				finalAnswer += "\n\n" + orchestrated.MainContent
			}
		} else {
			finalAnswer = orchestrated.MainContent
		}
		if finalAnswer == "" {
			finalAnswer = orchestrated.RawJSON // absolute last fallback
		}
		illustrationURL = orchestrated.IllustrationURL
		isStructuredJSON = orchestrated.IsStructuredJSON
		tokens = orchestrated.TokensUsed

		latency := int(time.Since(start).Milliseconds())
		s.log.Info("orchestrated response complete",
			zap.String("output_format", string(req.OutputFormat)),
			zap.Bool("is_structured_json", isStructuredJSON),
			zap.String("illustration_url", illustrationURL),
			zap.Int("orchestration_ms", orchestrated.OrchestrationMs),
			zap.Int("tokens_used", tokens),
		)
		go s.logInteraction(context.Background(), userID, req.Question, finalAnswer,
			string(req.OutputFormat), tokens, latency)

		return &AskResponse{
			Answer:           finalAnswer,
			OutputFormat:     req.OutputFormat,
			TokensUsed:       tokens,
			LatencyMs:        latency,
			ContextUsed:      ragCtx != "",
			ContextFound:     chunksFound,
			IllustrationURL:  illustrationURL,
			IsStructuredJSON: isStructuredJSON,
			// Typed sub-structs for frontend to render without parsing JSON
			Summary:  orchestrated.Summary,
			Detailed: orchestrated.Detailed,
			Anime:    orchestrated.Anime,
			Sports:   orchestrated.Sports,
			Academic: orchestrated.Academic,
		}, nil

	default:
		chatResp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
			Messages: []ChatMessage{
				{Role: "system", Content: sysPrompt},
				{Role: "user", Content: req.Question},
			},
			MaxTokens: 2048,
		})
		if err != nil {
			return nil, fmt.Errorf("qwen chat text failed: %w", err)
		}
		finalAnswer = chatResp.Content
		tokens = chatResp.TokensUsed
	}

	latency := int(time.Since(start).Milliseconds())

	s.log.Info("ask complete",
		zap.Int("total_prompt_tokens", tokens),
		zap.String("output_format", string(req.OutputFormat)),
	)

	go s.logInteraction(context.Background(), userID, req.Question, finalAnswer,
		string(req.OutputFormat), tokens, latency)

	return &AskResponse{
		Answer:           finalAnswer,
		OutputFormat:     req.OutputFormat,
		TokensUsed:       tokens,
		LatencyMs:        latency,
		ContextUsed:      ragCtx != "",
		ContextFound:     chunksFound,
		DownloadURL:      downloadURL,
		IllustrationURL:  illustrationURL,
		IsStructuredJSON: isStructuredJSON,
	}, nil
}

// Translate uses Qwen with a translation-specific system prompt.
func (s *AIService) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	sysPrompt := fmt.Sprintf(
		"You are a professional translator. Translate the following text accurately into %s. "+
			"Preserve the original meaning, tone, and formatting. "+
			"Return ONLY the translated text, nothing else.",
		req.TargetLang,
	)

	chatResp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: req.Text},
		},
		MaxTokens: 2000,
	})
	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}

	return &TranslateResponse{
		Translated: chatResp.Content,
		SourceLang: "auto-detect",
		TokensUsed: chatResp.TokensUsed,
	}, nil
}

// GeneratePPT generates slides from a topic/content string.
func (s *AIService) GeneratePPT(ctx context.Context, userID uuid.UUID, topic, content string, isRetry bool) (*PPTResult, error) {
	// Ask Qwen to generate slides JSON
	sysPrompt := `You are a presentation designer. Generate a structured slide deck in JSON format.
Return ONLY valid JSON with this exact structure (no markdown, no preamble):
{"slides":[{"title":"...","content":"...","speaker_notes":"..."}]}
Rules:
- 5-8 slides maximum
- Each title: max 60 characters
- Each content: bullet points separated by \n, max 200 characters total
- speaker_notes: conversational expansion, max 300 characters
- First slide: title slide with overview
- Last slide: summary or key takeaways`

	if isRetry {
		sysPrompt += "\n\nCRITICAL: Your previous response was invalid JSON. You MUST return ONLY valid JSON this time, no code blocks."
	}

	userMsg := fmt.Sprintf("Create a presentation about: %s\n\nContext:\n%s", topic, content)

	chatResp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: userMsg},
		},

		MaxTokens:   3000,
		Temperature: 0.4, // lower temp for structured JSON output
	})
	if err != nil {
		return nil, fmt.Errorf("generating slides content: %w", err)
	}

	return s.pptGen.Generate(ctx, PPTGenerateRequest{
		UserID:  userID.String(),
		Topic:   topic,
		Content: chatResp.Content,
	})
}

// GenerateAudio converts AI text to speech.
func (s *AIService) GenerateAudio(ctx context.Context, userID uuid.UUID, text, voice string) (*TTSResult, error) {
	return s.tts.GenerateAudio(ctx, userID.String(), TTSRequest{
		Text:  text,
		Voice: voice,
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// DB helpers
// ──────────────────────────────────────────────────────────────────────────────

type userProfile struct {
	learningStyle    string
	dominantInterest string
	explanationDepth string
}

func defaultProfile() userProfile {
	return userProfile{
		learningStyle:    "adaptive",
		dominantInterest: "general",
		explanationDepth: "intermediate",
	}
}

func (s *AIService) loadUserProfile(ctx context.Context, userID uuid.UUID) (userProfile, error) {
	const q = `
		SELECT
			COALESCE(ulp.preferences->>'learning_style', 'adaptive'),
			COALESCE(ulp.preferences->>'dominant_interest', 'general'),
			COALESCE(ulp.preferences->>'explanation_depth', 'intermediate')
		FROM user_learning_profiles ulp
		WHERE ulp.user_id = $1
		LIMIT 1`

	var p userProfile
	err := s.db.QueryRowContext(ctx, q, userID).Scan(
		&p.learningStyle,
		&p.dominantInterest,
		&p.explanationDepth,
	)
	if err != nil {
		return defaultProfile(), err
	}
	return p, nil
}

func (s *AIService) logInteraction(ctx context.Context, userID uuid.UUID,
	question, response, format string, tokens, latencyMs int) {

	summary := response
	if len(summary) > 500 {
		summary = summary[:497] + "..."
	}

	const q = `
		INSERT INTO ai_interactions
		  (user_id, question, output_format, response_summary, tokens_used, latency_ms, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,NOW())`

	if _, err := s.db.ExecContext(ctx, q,
		userID, question, format, summary, tokens, latencyMs); err != nil {
		s.log.Warn("failed to log ai interaction",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}
}

// Existing interface methods (backward compat with old provider pattern)
// These delegate to the new unified service.

func (s *AIService) Explain(ctx context.Context, req ExplainRequest) (*ExplainResponse, error) {
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		uid = uuid.Nil
	}

	interest := "general"
	if len(req.Interests) > 0 {
		interest = req.Interests[0]
	}

	format := OutputFormatDetailed
	if req.Style == "summary" {
		format = OutputFormatSummary
	}

	r, err := s.Ask(ctx, uid, AskRequest{
		Question:     fmt.Sprintf("Explain: %s (Subject: %s)", req.Topic, req.Subject),
		OutputFormat: format,
	})
	if err != nil {
		return nil, err
	}

	_ = interest
	return &ExplainResponse{
		Explanation: r.Answer,
		TokensUsed:  r.TokensUsed,
		LatencyMs:   r.LatencyMs,
	}, nil
}

func (s *AIService) GenerateIllustration(ctx context.Context, req IllustrationRequest) (*IllustrationResponse, error) {
	// Build a rich illustration description then ask Qwen to write vivid descriptive text.
	prompt := fmt.Sprintf(
		"Create a vivid, detailed visual description (not actual image) for an educational illustration about '%s'. "+
			"Style: %s. Themes from: %s. Description: %s. "+
			"Write as if describing the ideal educational infographic.",
		req.Topic, req.Style, strings.Join(req.Interests, ", "), req.Description,
	)

	resp, err := s.client.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 500,
	})
	if err != nil {
		return nil, err
	}

	return &IllustrationResponse{
		ImageURL:   fmt.Sprintf("/api/v1/ai/placeholder-illustration?topic=%s", req.Topic),
		Prompt:     resp.Content,
		TokensUsed: resp.TokensUsed,
	}, nil
}

func (s *AIService) AdaptStyle(ctx context.Context, req StyleRequest) (*StyleResponse, error) {
	return &StyleResponse{
		SystemPrompt: BuildSystemPrompt(PromptBuilderConfig{
			LearningStyle: req.PreferredStyle,
			DominantInterest: func() string {
				if len(req.Interests) > 0 {
					return req.Interests[0]
				}
				return "general"
			}(),
			ExplanationDepth: req.DifficultyLevel,
		}),
		Tone:     req.DifficultyLevel,
		Examples: strings.Join(req.Interests, ", "),
	}, nil
}
