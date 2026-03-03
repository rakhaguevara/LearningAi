package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/adaptive-ai-learn/backend/internal/ai"
	"github.com/adaptive-ai-learn/backend/internal/auth"
	"github.com/adaptive-ai-learn/backend/internal/config"
	"github.com/adaptive-ai-learn/backend/internal/learning"
	mw "github.com/adaptive-ai-learn/backend/internal/middleware"
	"github.com/adaptive-ai-learn/backend/internal/onboarding"
	"github.com/adaptive-ai-learn/backend/internal/personalization"
	"github.com/adaptive-ai-learn/backend/internal/user"
	jwtpkg "github.com/adaptive-ai-learn/backend/pkg/jwt"
)

func Setup(router *gin.Engine, db *sql.DB, cfg *config.Config, log *zap.Logger) {
	// ── Infrastructure ───────────────────────────────────
	jwtSvc := jwtpkg.NewService(cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshExpiry)

	// ── Auth ──────────────────────────────────────────────
	authRepo := auth.NewRepository(db, log)
	authSvc := auth.NewService(authRepo, jwtSvc, log)
	authHandler := auth.NewHandler(authSvc)

	// ── User ──────────────────────────────────────────────
	userRepo := user.NewRepository(db, log)
	userSvc := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userSvc)

	// ── AI Workspace ─────────────────────────────────────
	qwenClient := ai.NewQwenClient(cfg.Qwen.APIKey, cfg.Qwen.Endpoint, cfg.Qwen.Model, log)
	ragEngine := ai.NewRAGEngine(db, log)
	fileParser := ai.NewFileParser(qwenClient, log)
	pptGen := ai.NewPPTGenerator("", log)
	tts := ai.NewTTSService(cfg.Qwen.APIKey, cfg.Qwen.Endpoint, "", log)
	imageGen := ai.NewImageGenerator(cfg.Qwen.APIKey, cfg.Qwen.Endpoint, log)

	aiSvc := ai.NewAIService(qwenClient, ragEngine, fileParser, pptGen, tts, imageGen, db, log)
	aiHandler := ai.NewHandler(aiSvc, log)

	// ── Learning ──────────────────────────────────────────
	learningRepo := learning.NewRepository(db, log)
	learningSvc := learning.NewService(learningRepo, log)
	learningHandler := learning.NewHandler(learningSvc)

	// ── Personalization ──────────────────────────────────
	personalizationRepo := personalization.NewRepository(db, log)
	personalizationSvc := personalization.NewService(personalizationRepo, log)
	personalizationHandler := personalization.NewHandler(personalizationSvc)

	// ── Onboarding ────────────────────────────────────────
	onboardingRepo := onboarding.NewRepository(db, log)
	onboardingSvc := onboarding.NewService(onboardingRepo, log)
	onboardingHandler := onboarding.NewHandler(onboardingSvc)

	// ── Rate Limiter ─────────────────────────────────────
	rateLimiter := mw.NewRateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst)

	// ── Global Middleware ────────────────────────────────
	router.Use(mw.CORS(cfg.App.FrontendURL))
	router.Use(mw.RequestLogger(log))
	router.Use(mw.Recovery(log))
	router.Use(rateLimiter.Middleware())

	// ── Health ───────────────────────────────────────────
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "ailearn-api"})
	})

	// ── Auth Routes (public) ─────────────────────────────
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/google", authHandler.GoogleAuth)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// ── Protected Routes ─────────────────────────────────
	protected := router.Group("")
	protected.Use(mw.Auth(jwtSvc, log))
	{
		// User
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/profile", userHandler.GetProfile)
		}

		// ── AI Workspace ─────────────────────────────────
		aiGroup := protected.Group("/ai")
		{
			// Public format list (can be cached by frontend)
			aiGroup.GET("/formats", aiHandler.GetOutputFormats)

			// Core chat (rate-limited per endpoint)
			aiGroup.POST("/ask", aiHandler.Ask)

			// Document upload & RAG
			aiGroup.POST("/upload", aiHandler.UploadFile)
			aiGroup.GET("/sources", aiHandler.GetSources)

			// Generators
			aiGroup.POST("/generate-ppt", aiHandler.GeneratePPT)
			aiGroup.POST("/generate-audio", aiHandler.GenerateAudio)
			aiGroup.POST("/translate", aiHandler.Translate)

			// Downloads (user-scoped, path-safe)
			aiGroup.GET("/download/ppt/:user_id/:filename", aiHandler.DownloadPPT)
			aiGroup.GET("/download/audio/:user_id/:filename", aiHandler.DownloadAudio)

			// Legacy endpoints (backward compat)
			aiGroup.POST("/explain", aiHandler.Explain)
			aiGroup.POST("/generate-illustration", aiHandler.GenerateIllustration)
		}

		// Learning
		learningGroup := protected.Group("/learning")
		{
			learningGroup.POST("/start-session", learningHandler.StartSession)
		}

		// Personalization
		personalizationGroup := protected.Group("/personalization")
		{
			personalizationGroup.GET("/profile", personalizationHandler.GetProfile)
			personalizationGroup.GET("/prompt", personalizationHandler.GetAdaptivePrompt)
			personalizationGroup.GET("/learning-style", personalizationHandler.GetLearningStyle)
			personalizationGroup.GET("/interests", personalizationHandler.GetInterests)
			personalizationGroup.POST("/signal", personalizationHandler.RecordSignal)
			personalizationGroup.POST("/interest", personalizationHandler.AddInterest)
			personalizationGroup.POST("/feedback", personalizationHandler.RecordFeedback)
		}

		// Onboarding
		onboardingGroup := protected.Group("/onboarding")
		{
			onboardingGroup.GET("/status", onboardingHandler.GetStatus)
			onboardingGroup.POST("/submit", onboardingHandler.Submit)
		}

		// Profile updates
		profileGroup := protected.Group("/profile")
		{
			profileGroup.PUT("/update-learning", onboardingHandler.UpdateLearning)
		}
	}
}
