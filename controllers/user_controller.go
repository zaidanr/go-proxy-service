package controllers

import (
	"managed-proxy-server/helper"
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
	var user models.User
	var newUser models.NewUser

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		panic(err.Error())
	}

	if result := models.DB.Where("email = ?", newUser.Email).First(&user); result.Error != nil {
		// TODO: Differentiate between server error and user user not found error
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user = models.User{Username: newUser.Username, Email: newUser.Email, Hash: string(hash), UID: helper.GenerateUID()}
		if err := models.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, &user)
		return
	} else {
		c.JSON(http.StatusConflict, gin.H{"Error": "User already registered"})
		return
	}
}
