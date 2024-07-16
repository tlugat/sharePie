package refund

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve refunds: %v", err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Refund with ID %d not found", id)})
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
// @Success 200 {object} map[string.interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string.interface{} "Returns an error if the refund cannot be deleted"
// @Router /refunds/{id} [delete]
func (controller *Controller) DeleteRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.refundService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to delete refund with ID %d: %v", id, err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
