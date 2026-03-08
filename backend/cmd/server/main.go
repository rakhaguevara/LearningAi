package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/adaptive-ai-learn/backend/api/routes"
	"github.com/adaptive-ai-learn/backend/internal/common/logger"
	"github.com/adaptive-ai-learn/backend/internal/config"
	"github.com/adaptive-ai-learn/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.App.Env)
	defer log.Sync()

	log.Info("starting AI Learning Platform API",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// ── Database ────────────────────────────────────────
	db, err := database.Connect(cfg.DB, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := database.RunMigrations(db, log); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}

	// ── Router ──────────────────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// ✅ CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			cfg.App.FrontendURL,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(router, db, cfg, log)

	// ── Start Server ────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info("server listening", zap.String("addr", addr))

	if err := router.Run(addr); err != nil {
		log.Fatal("server failed to start", zap.Error(err))
	}
}
