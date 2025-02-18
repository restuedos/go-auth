package main

import (
	"go-auth/internal/config"
	"go-auth/internal/delivery/http/handler"
	"go-auth/internal/delivery/http/middleware"
	"go-auth/internal/delivery/http/routes"
	"go-auth/internal/domain/repository"
	"go-auth/internal/infrastructure/cache"
	"go-auth/internal/infrastructure/database"
	"go-auth/internal/usecase"
	"log"

	_ "go-auth/docs"

	"github.com/gin-gonic/gin"
	docs "github.com/swaggo/files"
	swag "github.com/swaggo/gin-swagger"
)

// @title Go Auth API
// @version 1.0
// @description This is a sample authentication API using Gin and Swaggo.
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer authentication
func main() {
	// Load environment variables
	cfg := config.LoadConfig()

	// Database connection
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Redis cache connection
	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisCache.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	database.SeedAdminUser(db, userRepo)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, redisCache)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)

	// Setup Gin router
	router := gin.Default()

	// Initialize middleware
	middleware.Init(cfg, authUsecase)

	// Setup Swagger
	router.GET("/swagger/*any", swag.WrapHandler(docs.Handler))

	// Setup routes
	routes.SetupRoutes(router, authHandler, userHandler)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
