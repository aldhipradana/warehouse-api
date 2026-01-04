package main

import (
	"github.com/aldhipradana/warehouse-api/middleware"
	"github.com/aldhipradana/warehouse-api/models"
	"github.com/aldhipradana/warehouse-api/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup DB
	db, _ := gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})

	r := gin.Default()

	// Use Action Logger Middleware
	r.Use(middleware.ActionLogger())

	// Register Routes
	routes.RegisterRoutes(r, db)

	r.Run(":8080")
}
