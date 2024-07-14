package refund

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/auth"
	"sharePie-api/internal/types"
	"strconv"
)

type Controller struct {
	refundService types.IRefundService
}

func NewController(service types.IRefundService) *Controller {
	return &Controller{refundService: service}
}

// FindRefunds retrieves all refunds.
// @Summary List all refunds
// @Description Retrieves a list of all refunds from the database
// @Tags Refunds
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of refunds"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /refunds [get]
func (controller *Controller) FindRefunds(c *gin.Context) {
	refunds, err := controller.refundService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": refunds})
}

// FindRefund retrieves a refund by ID.
// @Summary Get a single refund
// @Description Retrieves a refund by its ID from the database
// @Tags Refunds
// @Accept  json
// @Produce  json
// @Param id path int true "Refund ID"
// @Success 200 {object} map[string]interface{} "Returns the specified refund"
// @Failure 400 {object} map[string]interface{} "Returns an error if the refund is not found"
// @Router /refunds/{id} [get]
func (controller *Controller) FindRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	refund, err := controller.refundService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": refund})
}

// CreateRefund adds a new refund.
// @Summary Add a new refund
// @Description Adds a new refund to the database, linked to the authenticated user
// @Tags Refunds
// @Accept  json
// @Produce  json
// @Param input body types.CreateRefundInput true "Refund creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created refund"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /refunds [post]
func (controller *Controller) CreateRefund(c *gin.Context) {
	var input types.CreateRefundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := auth.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	refund, err := controller.refundService.Create(input, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": refund})
}

// UpdateRefund updates an existing refund.
// @Summary Update a refund
// @Description Updates an existing refund with new data
// @Tags Refunds
// @Accept  json
// @Produce  json
// @Param id path int true "Refund ID"
// @Param input body types.UpdateRefundInput true "Refund update data"
// @Success 200 {object} map[string]interface{} "Returns the updated refund"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the refund does not exist"
// @Router /refunds/{id} [patch]
func (controller *Controller) UpdateRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input types.UpdateRefundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refund, err := controller.refundService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": refund})
}

// DeleteRefund removes a refund.
// @Summary Delete a refund
// @Description Deletes a refund from the database
// @Tags Refunds
// @Accept  json
// @Produce  json
// @Param id path int true "Refund ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the refund cannot be deleted"
// @Router /refunds/{id} [delete]
func (controller *Controller) DeleteRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.refundService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete refund"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
