package routes

import (
	"github.com/aldhipradana/warehouse-api/middleware"
	"github.com/aldhipradana/warehouse-api/models"
	"github.com/aldhipradana/warehouse-api/restful"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUserRoutes sets up the routes for the User model
func RegisterUserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userCtrl := restful.NewCrudController[models.User](db)

	users := rg.Group("/users")
	users.Use(middleware.AuthMiddleware()) // All user routes require authentication
	{
		// Admin only routes
		users.GET("", middleware.AdminMiddleware(), userCtrl.Index)
		users.GET("/:id", middleware.AdminMiddleware(), userCtrl.Show)
		users.DELETE("/:id", middleware.AdminMiddleware(), userCtrl.Destroy)

		// Users can update themselves
		users.PUT("/:id", userCtrl.Update)
	}
}
