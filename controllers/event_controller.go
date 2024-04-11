package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/helpers"
	"go-project/services"
	"net/http"
	"strconv"
)

type EventController struct {
	eventService services.IEventService
}

func NewEventController(service services.IEventService) *EventController {
	return &EventController{eventService: service}
}

func (controller *EventController) FindEvents(c *gin.Context) {
	events, err := controller.eventService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func (controller *EventController) FindEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func (controller *EventController) CreateEvent(c *gin.Context) {
	var input services.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := helpers.GetUserFromContext(c)

	if !ok {
		return
	}

	event, err := controller.eventService.Create(input, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func (controller *EventController) UpdateEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := controller.eventService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func (controller *EventController) DeleteEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.eventService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
