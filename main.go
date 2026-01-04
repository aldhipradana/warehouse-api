package main

import (
	"log"

	"github.com/aldhipradana/warehouse-api/config"
	"github.com/aldhipradana/warehouse-api/middleware"
	"github.com/aldhipradana/warehouse-api/models"
	"github.com/aldhipradana/warehouse-api/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode based on debug setting
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup DB using config
	db, err := gorm.Open(sqlite.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.Product{}, &models.User{})

	// Initialize auth middleware with config
	middleware.InitAuth(cfg)

	r := gin.Default()

	// Use Action Logger Middleware
	r.Use(middleware.ActionLogger())

	// Register Routes
	routes.RegisterRoutes(r, db)

	// Run server with configured port
	r.Run(cfg.Server.GetServerAddress())
}
