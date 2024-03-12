package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sudhakarkandikattu/SplitWise/models"
)

func addGroup(context *gin.Context) {
	var group models.Group
	err := context.ShouldBindJSON(&group)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse requested data"})
		return
	}
	err = group.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the group"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "group created"})
}
func getUserGroups(context *gin.Context) {
	userId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "not a valid id for a group"})
	}
	groups, err := models.GetGroupsById(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the Groups by user id"})
		return
	}
	context.JSON(http.StatusOK, groups)
}
func getGroupById(context *gin.Context) {
	groupId, err := strconv.ParseInt(context.Param("groupid"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "not a valid id for a group"})
	}
	group, err := models.GetGroupByGroupId(groupId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the Groups by group id"})
		return
	}
	context.JSON(http.StatusOK, group)
}
