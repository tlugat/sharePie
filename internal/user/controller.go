package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/types"
	"strconv"
)

type Controller struct {
	userService types.IUserService
}

func NewController(service types.IUserService) *Controller {
	return &Controller{userService: service}
}

// FindUsers retrieves all users.
// @Summary List all users
// @Description Retrieves a list of all users from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of users"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /users [get]
func (controller *Controller) FindUsers(c *gin.Context) {
	users, err := controller.userService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// FindUser retrieves a user by ID.
// @Summary Get a single user
// @Description Retrieves a user by its ID from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Returns the specified user"
// @Failure 400 {object} map[string]interface{} "Returns an error if the user is not found"
// @Router /users/{id} [get]
func (controller *Controller) FindUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := controller.userService.FindOneById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser updates an existing user.
// @Summary Update a user
// @Description Updates an existing user with new data
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param input body services.UpdateUserInput true "User update data"
// @Success 200 {object} map[string]interface{} "Returns the updated user"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the user does not exist"
// @Router /users/{id} [put]
func (controller *Controller) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := controller.userService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser removes a user.
// @Summary Delete a user
// @Description Deletes a user from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the user cannot be deleted"
// @Router /users/{id} [delete]
func (controller *Controller) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (controller *Controller) GetUserFromToken(c *gin.Context) {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}
