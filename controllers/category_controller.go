package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/services"
	"net/http"
	"strconv"
)

type CategoryController struct {
	categoryService services.ICategoryService
}

func NewCategoryController(service services.ICategoryService) *CategoryController {
	return &CategoryController{categoryService: service}
}

func (controller *CategoryController) FindCategories(c *gin.Context) {
	categories, err := controller.categoryService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (controller *CategoryController) FindCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := controller.categoryService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (controller *CategoryController) CreateCategory(c *gin.Context) {
	var input services.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := controller.categoryService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (controller *CategoryController) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := controller.categoryService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (controller *CategoryController) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.categoryService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
