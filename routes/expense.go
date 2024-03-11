package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sudhakarkandikattu/SplitWise/models"
)

func addExpense(context *gin.Context) {
	var expense models.Expense
	err := context.ShouldBindJSON(&expense)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parse the Expense Data"})
		return
	}
	err = expense.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the expense"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Expense Added"})
}
