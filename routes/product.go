package routes

import (
	"github.com/aldhipradana/warehouse-api/middleware"
	"github.com/aldhipradana/warehouse-api/models"
	"github.com/aldhipradana/warehouse-api/restful"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterProductRoutes sets up the routes for the Product model
func RegisterProductRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	productCtrl := restful.NewCrudController[models.Product](db)

	products := rg.Group("/products")
	products.Use(middleware.AuthMiddleware()) // All routes require authentication
	{
		products.GET("", productCtrl.Index)
		products.GET("/:id", productCtrl.Show)
		products.POST("", productCtrl.Store)
		products.PUT("/:id", productCtrl.Update)
		products.DELETE("/:id", productCtrl.Destroy)
	}
}
