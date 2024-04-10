package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/services"
	"net/http"
	"strconv"
)

type UserController struct {
	userService services.IUserService
}

func NewUserController(service services.IUserService) *UserController {
	return &UserController{userService: service}
}

func (controller *UserController) FindUsers(c *gin.Context) {
	users, err := controller.userService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (controller *UserController) FindUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := controller.userService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (controller *UserController) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := controller.userService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (controller *UserController) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.userService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
