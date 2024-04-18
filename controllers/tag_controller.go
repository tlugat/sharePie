package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/services"
	"strconv"
)

type TagController struct {
	tagService services.ITagService
}

func NewTagController(service services.ITagService) *TagController {
	return &TagController{tagService: service}
}

// FindTags retrieves all tags.
// @Summary List all tags
// @Description Retrieves a list of all tags from the database
// @Tags Tags
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of tags"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /tags [get]
func (controller *TagController) FindTags(c *gin.Context) {
	tags, err := controller.tagService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// FindTag retrieves a tag by ID.
// @Summary Get a single tag
// @Description Retrieves a tag by its ID from the database
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200 {object} map[string]interface{} "Returns the specified tag"
// @Failure 400 {object} map[string]interface{} "Returns an error if the tag is not found"
// @Router /tags/{id} [get]
func (controller *TagController) FindTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := controller.tagService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// CreateTag adds a new tag.
// @Summary Add a new tag
// @Description Adds a new tag to the database
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param input body services.CreateTagInput true "Tag creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created tag"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid"
// @Router /tags [post]
func (controller *TagController) CreateTag(c *gin.Context) {
	var input services.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := controller.tagService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// UpdateTag updates an existing tag.
// @Summary Update a tag
// @Description Updates an existing tag with new data
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Param input body services.UpdateTagInput true "Tag update data"
// @Success 200 {object} map[string]interface{} "Returns the updated tag"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the tag does not exist"
// @Router /tags/{id} [put]
func (controller *TagController) UpdateTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := controller.tagService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// DeleteTag removes a tag.
// @Summary Delete a tag
// @Description Deletes a tag from the database
// @Tags Tags
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the tag cannot be deleted"
// @Router /tags/{id} [delete]
func (controller *TagController) DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.tagService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
