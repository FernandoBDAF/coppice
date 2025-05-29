package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fernandobarroso/profile-service/microservices/services/profile-service/internal/api/routes"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-service/internal/config"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-service/internal/domain/services"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-service/internal/pkg/logger"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-service/internal/session"
	"go.uber.org/zap"
)

func main() {
	// Create background context
	ctx := context.Background()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	loggerCfg := &logger.Config{
		Level:       cfg.Logging.Level,
		Environment: cfg.Environment,
		ServiceName: "profile-api",
		Format:      cfg.Logging.Format,
		LogFile:     cfg.Logging.LogFile,
		Shipping: &logger.ShippingConfig{
			Enabled:    cfg.Logging.Shipping.Enabled,
			Endpoint:   cfg.Logging.Shipping.Endpoint,
			BufferSize: cfg.Logging.Shipping.BufferSize,
			MaxRetries: cfg.Logging.Shipping.MaxRetries,
			RetryDelay: cfg.Logging.Shipping.RetryDelay,
		},
	}
	if err := logger.Initialize(loggerCfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize auth service client
	authClient := services.NewAuthServiceClient(cfg)

	// Initialize session manager
	logger.LogInfo(ctx, "Initializing Redis session manager")
	sessionManager, err := session.NewSessionManager(authClient)
	if err != nil {
		logger.LogError(ctx, "Failed to initialize session manager", err)
		log.Fatalf("Failed to initialize session manager: %v", err)
	}
	defer sessionManager.Close()

	// Initialize router
	router := api.NewRouter(cfg, authClient, sessionManager)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.LogInfo(ctx, "Starting server",
		zap.String("address", addr))
	if err := router.Run(addr); err != nil {
		logger.LogError(ctx, "Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
