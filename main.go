package main

import (
	"github.com/aldhipradana/warehouse-api/models"
	"github.com/aldhipradana/warehouse-api/restful"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup DB
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})

	r := gin.Default()

	// Initialize the Generic Controller for "Product"
	productCtrl := restful.NewCrudController[models.Product](db)

	// 4. Register Routes
	api := r.Group("/api/products")
	{
		api.GET("", productCtrl.Index) // Supports filter={"price": {"operator":">", "value":100}}
		api.GET("/:id", productCtrl.Show)
		api.POST("", productCtrl.Store)
		api.PUT("/:id", productCtrl.Update)
		api.DELETE("/:id", productCtrl.Destroy)
	}

	r.Run(":8080")
}
