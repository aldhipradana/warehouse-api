// restful/controller.go
package restful

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationResponse matches Laravel's pagination structure
type PaginationResponse[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// CrudController is a generic controller for any model T
type CrudController[T any] struct {
	DB *gorm.DB
}

// NewCrudController creates a new instance
func NewCrudController[T any](db *gorm.DB) *CrudController[T] {
	return &CrudController[T]{DB: db}
}

// Index - GET /api/resource
func (c *CrudController[T]) Index(ctx *gin.Context) {
	var items []T
	var total int64
	var model T // Zero value of T just to access methods/traits

	// 1. Get Query Params
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	sort := ctx.DefaultQuery("sort", "created_at")
	order := ctx.DefaultQuery("order", "desc")
	filterJSON := ctx.Query("filter")
	relations := ctx.Query("relations")

	// search can be in "q" or inside filter json, handling "q" separately for ease
	search := ctx.Query("q")

	// 2. Start Query
	query := c.DB.Model(&model)

	// 3. Eager Load Relations
	if relations != "" {
		for _, rel := range strings.Split(relations, ",") {
			query = query.Preload(strings.TrimSpace(rel))
		}
	}

	// 4. Apply Filters (The Trait Logic)
	query = ApplyFilters(query, filterJSON, search, model)

	// 5. Sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// 6. Pagination count
	query.Count(&total)

	// 7. Fetch Data
	offset := (page - 1) * limit
	result := query.Limit(limit).Offset(offset).Find(&items)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 8. Return Response
	ctx.JSON(http.StatusOK, PaginationResponse[T]{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// Show - GET /api/resource/:id
func (c *CrudController[T]) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	var item T

	// Handle Relations
	query := c.DB
	if relations := ctx.Query("relations"); relations != "" {
		for _, rel := range strings.Split(relations, ",") {
			query = query.Preload(strings.TrimSpace(rel))
		}
	}

	if err := query.First(&item, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// Store - POST /api/resource
func (c *CrudController[T]) Store(ctx *gin.Context) {
	var item T
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// Update - PUT /api/resource/:id
func (c *CrudController[T]) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var item T

	// Check existence
	if err := c.DB.First(&item, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Bind new data
	// Note: BindJSON will overwrite fields in 'item'.
	// In Go, partial updates usually require a map[string]interface{}
	// but for simplicity we bind to struct here.
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save
	if err := c.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// Destroy - DELETE /api/resource/:id
func (c *CrudController[T]) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")
	var item T

	if err := c.DB.Delete(&item, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}
