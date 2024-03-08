package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sudhakarkandikattu/SplitWise/db"
	"github.com/sudhakarkandikattu/SplitWise/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)
	server.Run(":8080")
}

func registerUser(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Hello There"})
}
