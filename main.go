package main

import (
	"managed-proxy-server/controllers"
	"managed-proxy-server/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectDataBase()

	r.GET("/users", controllers.GetUser)
	r.POST("/register", controllers.AddUser)

	r.Run("localhost:9090")
}
