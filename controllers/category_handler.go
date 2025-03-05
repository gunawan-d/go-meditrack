package controllers

import (
	"meditrack/repository"
	"meditrack/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Repo *repository.CategoryRepository
}

func NewCategoryController(repo *repository.CategoryRepository) *CategoryController {
	return &CategoryController{Repo: repo}
}

// Create Category
func (c *CategoryController) CreateCategoryHandler(ctx *gin.Context) {
	var category structs.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id, err := c.Repo.CreateCategory(category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category_id": id})
}

// Get All Categories
func (c *CategoryController) GetAllCategoriesHandler(ctx *gin.Context) {
	categories, err := c.Repo.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}
