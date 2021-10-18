package controllers

import (
	"managed-proxy-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *gin.Context) {
	var user models.User
	if result := models.DB.First(&user); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Error": result.Error.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, result)

	}
	return
}

func AddUser(c *gin.Context) {
	var newUser models.NewUser
	var uid string = "123"

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		panic(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user := models.User{Username: newUser.Username, Email: newUser.Email, Hash: string(hash), UID: string(uid)}
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
