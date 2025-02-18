package routes

import (
	"go-auth/internal/delivery/http/handler"
	"go-auth/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	api := router.Group("/api")
	{
		// Public routes - Authentication related
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)

		// Protected routes - Admin related
		admin := api.Group("/admin", middleware.AuthMiddleware())
		{
			admin.GET("/users", userHandler.GetAllUsers)
		}

		// Protected routes - User related
		user := api.Group("/user", middleware.AuthMiddleware())
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.DELETE("/profile", userHandler.DeleteProfile)
		}
	}
}
