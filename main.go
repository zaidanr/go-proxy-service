package main

import (
	"log"
	"managed-proxy-server/auth"
	"managed-proxy-server/middleware"
	"managed-proxy-server/models"
	"managed-proxy-server/user"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
)

func main() {
	identityKey := "id"
	r := gin.Default()
	models.ConnectDataBase()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "go-proxy-service",
		Key:           []byte(os.Getenv("JWTSECRET")),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   identityKey,
		PayloadFunc:   auth.MapClaims,
		Authenticator: auth.Login,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// r.GET("/users", user.GetUser)
	// r.POST("/register", controllers.AddUser)
	r.POST("/login", authMiddleware.LoginHandler)

	admin := r.Group("/admin")

	admin.Use(authMiddleware.MiddlewareFunc())
	{
		admin.GET("/users", middleware.Authorize("admin_action", "read", fileadapter.NewAdapter("config/policy.csv")), user.GetUser)
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims["id"])
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.Run("localhost:9090")
}
