package main

import (
	"log"

	"github.com/aldhipradana/warehouse-api/database/seed"
	"github.com/aldhipradana/warehouse-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup DB (using the same path as main app)
	db, err := gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Ensure tables exist
	log.Println("Migrating database...")
	db.AutoMigrate(&models.Product{})

	// Run Seeder
	log.Println("Seeding products...")
	if err := seed.SeedProducts(db); err != nil {
		log.Fatal("failed to seed products:", err)
	}

	log.Println("Seeding completed successfully!")
}
