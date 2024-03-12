package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sudhakarkandikattu/SplitWise/models"
)

func getUserToUserOwes(context *gin.Context) {
	groupId, err := strconv.ParseInt(context.Param("groupid"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "not a valid id for a group"})
	}
	userId, err := strconv.ParseInt(context.Param("userid"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "not a valid id for a user"})
	}
	userToUserOwes, err := models.GetUserToUserOwesByGroupId(groupId, userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the UserToUserOwes by group id"})
	}
	context.JSON(http.StatusOK, userToUserOwes)
}
