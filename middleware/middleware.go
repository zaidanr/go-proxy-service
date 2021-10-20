package middleware

import (
	"fmt"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/gin-gonic/gin"
)

func Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// err := auth.TokenValid(c.Request)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, "user hasn't logged in yet")
		// 	c.Abort()
		// 	return
		// }
		// metadata, err := auth.ExtractTokenMetadata(c.Request)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, "unauthorized")
		// 	return
		// }

		// casbin enforces policy
		claims := jwt.ExtractClaims(c)
		role := claims["role"].(string)
		ok, err := enforce(role, obj, act, adapter)
		//ok, err := enforce(val.(string), obj, act, adapter)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok := enforcer.Enforce(sub, obj, act)
	return ok, nil
}
