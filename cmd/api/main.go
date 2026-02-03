package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	_ "github.com/lib/pq"

	"github.com/yourusername/resume-builder/internal/client"
	"github.com/yourusername/resume-builder/internal/config"
	"github.com/yourusername/resume-builder/internal/crypto"
	"github.com/yourusername/resume-builder/internal/handler"
	"github.com/yourusername/resume-builder/internal/repository"
	"github.com/yourusername/resume-builder/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Info("starting application", "env", cfg.Server.Env)

	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("database connected")

	// Initialize Redis cache
	cache := client.NewCacheClient(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.Enabled,
	)
	defer cache.Close()
	if cfg.Redis.Enabled {
		logger.Info("redis cache enabled", "addr", cfg.Redis.Addr)
	}

	encryptor, err := crypto.NewEncryptor(cfg.Crypto.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to create encryptor: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resumeRepo := repository.NewResumeRepository(db)

	// Initialize clients
	githubClient := client.NewGitHubClient()
	llmClient := client.NewLLMClient(cfg.OpenAI.APIKey, cfg.OpenAI.Enabled)
	if cfg.OpenAI.Enabled {
		logger.Info("llm enabled for resume summaries")
	}

	// Initialize services
	jwtService := service.NewJWTService(cfg.Crypto.JWTSecret)
	authService := service.NewAuthService(
		userRepo,
		githubClient,
		encryptor,
		cfg.GitHub.ClientID,
		cfg.GitHub.ClientSecret,
		cfg.GitHub.RedirectURL,
	)
	githubService := service.NewGitHubService(githubClient, cache)
	rankingService := service.NewRankingService()
	resumeService := service.NewResumeService(resumeRepo, userRepo, githubService, rankingService, llmClient)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, jwtService)
	resumeHandler := handler.NewResumeHandler(resumeService, authService)
	authMiddleware := handler.NewAuthMiddleware(jwtService, logger)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(httprate.LimitByIP(50, 1*time.Minute))

		r.Post("/resumes/generate", resumeHandler.Generate)
		r.Get("/resumes", resumeHandler.List)
		r.Get("/resumes/{id}", resumeHandler.Get)
		r.Put("/resumes/{id}", resumeHandler.Update)
		r.Delete("/resumes/{id}", resumeHandler.Delete)
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("server starting", "port", cfg.Server.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		logger.Info("shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		logger.Info("server stopped gracefully")
	}

	return nil
}
