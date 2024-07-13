package avatar

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/types"
	"strconv"
)

type Controller struct {
	avatarService types.IAvatarService
}

func NewController(service types.IAvatarService) *Controller {
	return &Controller{avatarService: service}
}

// FindAvatars retrieves all avatars.
// @Summary List all avatars
// @Description Retrieves a list of all avatars from the database
// @Tags Avatars
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of avatars"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /avatars [get]
func (controller *Controller) FindAvatars(c *gin.Context) {
	avatars, err := controller.avatarService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": avatars})
}

// FindAvatar retrieves an avatar by ID.
// @Summary Get a single avatar
// @Description Retrieves an avatar by its ID from the database
// @Tags Avatars
// @Accept  json
// @Produce  json
// @Param id path int true "Avatar ID"
// @Success 200 {object} map[string]interface{} "Returns the specified avatar"
// @Failure 400 {object} map[string]interface{} "Returns an error if the avatar is not found"
// @Router /avatars/{id} [get]
func (controller *Controller) FindAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	avatar, err := controller.avatarService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": avatar})
}

// CreateAvatar adds a new avatar.
// @Summary Add a new avatar
// @Description Adds a new avatar to the database, linked to the authenticated user
// @Tags Avatars
// @Accept  json
// @Produce  json
// @Param input body types.CreateAvatarInput true "Avatar creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created avatar"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /avatars [post]
func (controller *Controller) CreateAvatar(c *gin.Context) {
	var input types.CreateAvatarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	avatar, err := controller.avatarService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": avatar})
}

// UpdateAvatar updates an existing avatar.
// @Summary Update an avatar
// @Description Updates an existing avatar with new data
// @Tags Avatars
// @Accept  json
// @Produce  json
// @Param id path int true "Avatar ID"
// @Param input body types.UpdateAvatarInput true "Avatar update data"
// @Success 200 {object} map[string]interface{} "Returns the updated avatar"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the avatar does not exist"
// @Router /avatars/{id} [patch]
func (controller *Controller) UpdateAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateAvatarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	avatar, err := controller.avatarService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": avatar})
}

// DeleteAvatar removes an avatar.
// @Summary Delete an avatar
// @Description Deletes an avatar from the database
// @Tags Avatars
// @Accept  json
// @Produce  json
// @Param id path int true "Avatar ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the avatar cannot be deleted"
// @Router /avatars/{id} [delete]
func (controller *Controller) DeleteAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.avatarService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete avatar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
