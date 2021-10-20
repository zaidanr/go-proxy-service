package auth

import (
	"log"
	"managed-proxy-server/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) (interface{}, error) {
	var loginVals models.LoginUser
	var user models.User
	if err := c.ShouldBindJSON(&loginVals); err != nil {
		log.Fatal(err)
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	if result := models.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return "", jwt.ErrFailedAuthentication
	} else {
		res := result.Value.(*models.User)
		if err := bcrypt.CompareHashAndPassword([]byte(res.Hash), []byte(loginVals.Password)); err != nil {
			return "", jwt.ErrFailedAuthentication
		}
		return res, nil
	}
}

func MapClaims(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"id":   v.Email,
			"role": v.Role,
		}
	}
	return jwt.MapClaims{}
}
