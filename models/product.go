package models

import "gorm.io/gorm"

// Product represents a product in the system
type Product struct {
	gorm.Model
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

// GetSearchableFields returns the fields that can be searched/filtered
func (Product) GetSearchableFields() []string {
	return []string{"name", "status"}
}
