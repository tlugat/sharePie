package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/services"
	"net/http"
	"strconv"
)

type TagController struct {
	tagService services.ITagService
}

func NewTagController(service services.ITagService) *TagController {
	return &TagController{tagService: service}
}

func (controller *TagController) FindTags(c *gin.Context) {
	tags, err := controller.tagService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func (controller *TagController) FindTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := controller.tagService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tag})
}

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

func (controller *TagController) UpdateTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := controller.tagService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tag})
}

func (controller *TagController) DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.tagService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
