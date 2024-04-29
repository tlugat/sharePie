package expense

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sharePie-api/internal/auth"
	"strconv"
)

type Controller struct {
	expenseService IExpenseService
}

func NewController(service IExpenseService) *Controller {
	return &Controller{expenseService: service}
}

// FindExpenses retrieves all expenses.
// @Summary List all expenses
// @Description Retrieves a list of all expenses from the database
// @Tags Expenses
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns a list of expenses"
// @Failure 500 {object} map[string]interface{} "Returns an error if the request fails"
// @Router /expenses [get]
func (controller *Controller) FindExpenses(c *gin.Context) {
	expenses, err := controller.expenseService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

// FindExpense retrieves an expense by ID.
// @Summary Get a single expense
// @Description Retrieves an expense by its ID from the database
// @Tags Expenses
// @Accept  json
// @Produce  json
// @Param id path int true "Expense ID"
// @Success 200 {object} map[string]interface{} "Returns the specified expense"
// @Failure 400 {object} map[string]interface{} "Returns an error if the expense is not found"
// @Router /expenses/{id} [get]
func (controller *Controller) FindExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	expense, err := controller.expenseService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

// CreateExpense adds a new expense.
// @Summary Add a new expense
// @Description Adds a new expense to the database, linked to the authenticated user
// @Tags Expenses
// @Accept  json
// @Produce  json
// @Param input body services.CreateExpenseInput true "Expense creation data"
// @Success 200 {object} map[string]interface{} "Returns the newly created expense"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or user authentication fails"
// @Router /expenses [post]
func (controller *Controller) CreateExpense(c *gin.Context) {
	var input CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := auth.GetUserFromContext(c)

	if !ok {
		return
	}

	expense, err := controller.expenseService.Create(input, *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

// UpdateExpense updates an existing expense.
// @Summary Update an expense
// @Description Updates an existing expense with new data
// @Tags Expenses
// @Accept  json
// @Produce  json
// @Param id path int true "Expense ID"
// @Param input body services.UpdateExpenseInput true "Expense update data"
// @Success 200 {object} map[string]interface{} "Returns the updated expense"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid or the expense does not exist"
// @Router /expenses/{id} [put]
func (controller *Controller) UpdateExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input UpdateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	expense, err := controller.expenseService.Update(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

// DeleteExpense removes an expense.
// @Summary Delete an expense
// @Description Deletes an expense from the database
// @Tags Expenses
// @Accept  json
// @Produce  json
// @Param id path int true "Expense ID"
// @Success 200 {object} map[string]interface{} "Confirms successful deletion"
// @Failure 400 {object} map[string]interface{} "Returns an error if the expense cannot be deleted"
// @Router /expenses/{id} [delete]
func (controller *Controller) DeleteExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.expenseService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
