package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	//posts
	server.POST("/register", register)
	server.POST("/addgroup", addGroup)
	server.POST("/addexpense", addExpense)
	//gets
	server.GET("/getusergroups/:id", getUserGroups)
	server.GET("/allusers", getUsers)
	server.GET("/group/:groupid/user/:userid", getUserToUserOwes)
	server.GET("/getgroup/:groupid", getGroupById)

}
