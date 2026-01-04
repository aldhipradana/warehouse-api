package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up all the API routes for the application
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		RegisterProductRoutes(api, db)
		// RegisterCategoryRoutes(api, db) // Example of adding another model
	}
}
