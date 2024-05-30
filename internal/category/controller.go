package category

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/types"
	"strconv"
)

type Controller struct {
	categoryService types.ICategoryService
}

func NewController(service types.ICategoryService) *Controller {
	return &Controller{categoryService: service}
}

// FindCategories retrieves all categories.
// @Summary List all categories
// @Description Retrieves a list of all categories from the database
// @Tags Categories
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of categories"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /categories [get]
func (controller *Controller) FindCategories(c *gin.Context) {
	categories, err := controller.categoryService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// FindCategory retrieves a category by ID.
// @Summary Get a single category
// @Description Retrieves a category by its ID from the database
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Returns the specified category"
// @Failure 400 {object} map[string]interface{} "Returns an error if the category is not found"
// @Router /categories/{id} [get]
func (controller *Controller) FindCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := controller.categoryService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory adds a new category.
// @Summary Add a new category
// @Description Adds a new category to the database
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param input body services.CreateCategoryInput true "Category creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created category"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid"
// @Router /categories [post]
func (controller *Controller) CreateCategory(c *gin.Context) {
	var input types.CreateCategoryInput
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

// UpdateCategory updates an existing category.
// @Summary Update a category
// @Description Updates an existing category with new data
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Param input body services.UpdateCategoryInput true "Category update data"
// @Success 200 {object} map[string]interface{} "Returns the updated category"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the category does not exist"
// @Router /categories/{id} [put]
func (controller *Controller) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := controller.categoryService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory removes a category.
// @Summary Delete a category
// @Description Deletes a category from the database
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the category cannot be deleted"
// @Router /categories/{id} [delete]
func (controller *Controller) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.categoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
