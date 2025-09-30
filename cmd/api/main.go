package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet/internal/config"
	"wallet/internal/domain"
	"wallet/internal/handler"
	"wallet/internal/infrastructure/cache"
	postgresRepo "wallet/internal/infrastructure/postgres"
	"wallet/internal/infrastructure/redis"
	"wallet/internal/usecase"

	"wallet/docs" // Import the generated docs

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Wallet App API
// @version 1.0
// @description This is a sample wallet API.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Load Configuration using Viper
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	// 2. Initialize Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 3. Initialize Sentry
	if err := sentry.Init(sentry.ClientOptions{Dsn: cfg.SentryDSN}); err != nil {
		slog.Error("Sentry initialization failed", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	// 4. Connect to Database
	db, err := gorm.Open(postgres.Open(cfg.DBSource), &gorm.Config{})
	if err != nil {
		slog.Error("Cannot connect to database", "error", err)
		sentry.CaptureException(err)
		os.Exit(1)
	}
	db.AutoMigrate(&domain.User{}, &domain.Wallet{})

	// 5. Dependency Injection (Wiring)
	postgresUserRepo := postgresRepo.NewPostgresUserRepository(db)
	cacheRepo, err := redis.NewRedisCacheRepository(cfg.RedisAddr)
	if err != nil {
		slog.Error("Cannot connect to Redis", "error", err)
		sentry.CaptureException(err)
		os.Exit(1)
	}
	// Wrap the postgres repo with the cache decorator
	userRepo := cache.NewCachedUserRepository(cacheRepo, postgresUserRepo)

	walletRepo := postgresRepo.NewPostgresWalletRepository(db)
	txnRepo := postgresRepo.NewPostgresTxnRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo, walletRepo, txnRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, txnRepo, logger)

	userHandler := handler.NewUserHandler(userUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase, logger)

	// 6. Setup Web Server (Fiber)
	app := fiber.New()

	// Swagger documentation endpoints
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		return c.JSON(docs.SwaggerInfo)
	})

	// Simple Swagger UI redirect
	app.Get("/swagger", func(c fiber.Ctx) error {
		swaggerURL := "https://petstore.swagger.io/?url=" + c.BaseURL() + "/swagger/doc.json"
		return c.Redirect().To(swaggerURL)
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/users", userHandler.CreateUser)
	v1.Post("/wallets/recharge", walletHandler.Recharge)
	v1.Post("/wallets/transfer", walletHandler.Transfer)

	// 7. Start Server with Graceful Shutdown
	port := cfg.ServerPort
	if port == "" {
		port = "8080" // Default port
	}

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server", "port", port)
		if err := app.Listen(":" + port); err != nil {
			slog.Error("Failed to start server", "error", err)
			sentry.CaptureException(err)
		}
	}()

	// Wait for interrupt signal
	<-c
	slog.Info("Shutting down server gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		sentry.CaptureException(err)
	}

	// Close database connections
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}

	slog.Info("Server shutdown complete")
}
