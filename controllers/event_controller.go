package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/helpers"
	"sharePie-api/services"
	"strconv"
)

type EventController struct {
	eventService   services.IEventService
	balanceService services.IEventBalanceService
}

func NewEventController(service services.IEventService, balanceService services.IEventBalanceService) *EventController {
	return &EventController{eventService: service, balanceService: balanceService}
}

// FindEvents retrieves all events.
// @Summary List all events
// @Description Retrieves a list of all events from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of events"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /events [get]
func (controller *EventController) FindEvents(c *gin.Context) {
	events, err := controller.eventService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": events})
}

// FindEvent retrieves an event by ID.
// @Summary Get a single event
// @Description Retrieves an event by its ID from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns the specified event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event is not found"
// @Router /events/{id} [get]
func (controller *EventController) FindEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

// CreateEvent adds a new event.
// @Summary Add a new event
// @Description Adds a new event to the database, linked to the authenticated user
// @Tags Events
// @Accept  json
// @Produce  json
// @Param input body services.CreateEventInput true "Event creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /events [post]
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

// UpdateEvent updates an existing event.
// @Summary Update an event
// @Description Updates an existing event with new data
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param input body services.UpdateEventInput true "Event update data"
// @Success 200 {object} map[string]interface{} "Returns the updated event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the event does not exist"
// @Router /events/{id} [put]
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

// DeleteEvent removes an event.
// @Summary Delete an event
// @Description Deletes an event from the database
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event cannot be deleted"
// @Router /events/{id} [delete]
func (controller *EventController) DeleteEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.eventService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// GetEventUsers retrieves all users for an event.
// @Summary Get event users
// @Description Retrieves all users for a specified event
// @Tags Events
// @Accept  json
// @Produce  json
// @Param eventId path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns a list of users for the event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event does not exist"
// @Router /events/{id}/users [get]
func (controller *EventController) GetEventUsers(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	users, err := controller.eventService.GetUsers(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetEventBalanceSummary retrieves a summary of balances for an event.
// @Summary Get event balance summary
// @Description Retrieves a summary of balances for a specified event
// @Tags Events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Returns a list of balances for the event"
// @Failure 400 {object} map[string]interface{} "Returns an error if the event does not exist"
// @Router /events/{id}/summary [get]
func (controller *EventController) GetEventBalanceSummary(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	event, err := controller.eventService.FindOne(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found!"})
		return
	}

	balanceSummary, err := controller.balanceService.GetBalanceSummary(event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": balanceSummary})
}
