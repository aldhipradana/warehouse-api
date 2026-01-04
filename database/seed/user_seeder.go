package seed

import (
	"github.com/aldhipradana/warehouse-api/models"
	"gorm.io/gorm"
)

// SeedUsers populates the database with initial user data
func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{Name: "Admin User", Email: "admin@example.com", Password: "admin123", Role: "admin"},
		{Name: "John Doe", Email: "john@example.com", Password: "password123", Role: "user"},
		{Name: "Jane Smith", Email: "jane@example.com", Password: "password123", Role: "user"},
		{Name: "Bob Manager", Email: "bob@example.com", Password: "password123", Role: "manager"},
	}

	for _, u := range users {
		// Check if user already exists to avoid duplicates
		var count int64
		db.Model(&models.User{}).Where("email = ?", u.Email).Count(&count)
		if count == 0 {
			if err := db.Create(&u).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
