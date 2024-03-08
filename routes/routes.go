package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", register)
	server.GET("/allusers", getUsers)
}
