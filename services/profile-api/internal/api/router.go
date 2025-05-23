package api

import (
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/config"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/handlers"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/middleware"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/services"
	"github.com/gin-gonic/gin"
)

// Router represents the API router
type Router struct {
	engine         *gin.Engine
	cfg            *config.Config
	authClient     *services.AuthServiceClient
	sessionManager handlers.SessionManagerInterface
}

// NewRouter creates a new API router
func NewRouter(cfg *config.Config, authClient *services.AuthServiceClient, sessionManager handlers.SessionManagerInterface) *Router {
	router := &Router{
		engine:         gin.Default(),
		cfg:            cfg,
		authClient:     authClient,
		sessionManager: sessionManager,
	}
	router.setupRoutes()
	return router
}

// setupRoutes configures all the routes for the API
func (r *Router) setupRoutes() {
	// Health check endpoint
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Metrics endpoints
	metrics := r.engine.Group("/metrics")
	{
		metricsHandler := handlers.NewMetricsHandler()
		metrics.GET("", metricsHandler.GetMetrics)
		metrics.DELETE("", metricsHandler.ResetMetrics)
	}

	// API v1 group
	v1 := r.engine.Group("/api/v1")
	{
		// Auth endpoints
		auth := v1.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler(r.sessionManager)
			auth.POST("/token", authHandler.Authenticate)
			auth.POST("/validate", authHandler.ValidateToken)
		}

		// Profile endpoints
		profiles := v1.Group("/profiles")
		{
			// Initialize storage client
			storageClient := services.NewStorageClient(r.cfg)

			// Initialize profile service
			profileService := services.NewProfileService(r.cfg, storageClient)

			// Initialize profile handler
			profileHandler := handlers.NewProfileHandler(profileService)

			// Apply session middleware
			profiles.Use(middleware.SessionMiddleware(r.sessionManager))

			// Profile routes
			profiles.GET("", profileHandler.GetProfiles)
			profiles.GET("/:id", profileHandler.GetProfile)
			profiles.POST("", profileHandler.CreateProfile)
			profiles.PUT("/:id", profileHandler.UpdateProfile)
			profiles.DELETE("/:id", profileHandler.DeleteProfile)
		}
	}
}

// Run starts the API server
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
