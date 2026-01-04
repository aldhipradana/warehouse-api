package seed

import (
	"github.com/aldhipradana/warehouse-api/models"
	"gorm.io/gorm"
)

// SeedProducts populates the database with initial product data
func SeedProducts(db *gorm.DB) error {
	products := []models.Product{
		{Name: "Laptop Pro", Price: 1500.00, Status: "active"},
		{Name: "Wireless Mouse", Price: 25.50, Status: "active"},
		{Name: "Mechanical Keyboard", Price: 89.99, Status: "active"},
		{Name: "USB-C Hub", Price: 45.00, Status: "inactive"},
		{Name: "Monitor 4K", Price: 350.00, Status: "active"},
		{Name: "Gaming Chair", Price: 299.99, Status: "active"},
		{Name: "Webcam HD", Price: 59.99, Status: "active"},
		{Name: "Desk Lamp", Price: 19.99, Status: "inactive"},
		{Name: "External SSD 1TB", Price: 120.00, Status: "active"},
		{Name: "Noise Cancelling Headphones", Price: 199.00, Status: "active"},
	}

	for _, p := range products {
		// Check if product already exists to avoid duplicates
		var count int64
		db.Model(&models.Product{}).Where("name = ?", p.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&p).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
