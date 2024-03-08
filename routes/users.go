package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sudhakarkandikattu/SplitWise/models"
)

func register(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse requested data"})
		return
	}
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "user created"})
}
func getUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the users"})
		return
	}
	context.JSON(http.StatusOK, users)
}
