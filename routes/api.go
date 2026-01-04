package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up all the API routes for the application
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		// Authentication routes (public)
		RegisterAuthRoutes(api, db)

		// User management routes (protected)
		RegisterUserRoutes(api, db)

		// Product routes (all protected)
		RegisterProductRoutes(api, db)
	}
}
