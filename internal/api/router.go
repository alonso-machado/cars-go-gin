package api

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/username/go-car-service/internal/repository"
	"github.com/username/go-car-service/internal/service"
	"github.com/username/go-car-service/pkg/logger"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(engine *gin.Engine, db *sql.DB) {
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	engine.Use(cors.New(config))

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	apiV1 := engine.Group("/api/v1")


	// Initialize repositories
	carRepo := repository.NewCarRepository(db)

	// Initialize services
	carService := service.NewCarService(carRepo)

	// Initialize handlers
	carHandler := NewCarHandler(carService)

	// Register routes
	carHandler.RegisterRoutes(apiV1)


	// 404 handler
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Endpoint not found",
		})
	})

	// Log all requests
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logger.GetLogger().Writer(),
	}))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Errorf("Panic recovered: %s", err)
			c.JSON(500, ErrorResponse{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		c.AbortWithStatus(500)
	}))
}
