package controllers

import (
	"github.com/gin-gonic/gin"
	"go-project/helpers"
	"go-project/services"
	"net/http"
	"strconv"
)

type ExpenseController struct {
	expenseService services.IExpenseService
}

func NewExpenseController(service services.IExpenseService) *ExpenseController {
	return &ExpenseController{expenseService: service}
}

func (controller *ExpenseController) FindExpenses(c *gin.Context) {
	expenses, err := controller.expenseService.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

func (controller *ExpenseController) FindExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	expense, err := controller.expenseService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expense})
}

func (controller *ExpenseController) CreateExpense(c *gin.Context) {
	var input services.CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := helpers.GetUserFromContext(c)

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

func (controller *ExpenseController) UpdateExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input services.UpdateExpenseInput
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

func (controller *ExpenseController) DeleteExpense(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := controller.expenseService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete expense"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
