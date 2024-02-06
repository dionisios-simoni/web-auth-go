package main

import (
	"github.com/gin-gonic/gin"
	"github.com/web-auth-go/03_jwt/handlers"
	"github.com/web-auth-go/03_jwt/initialisers"
	"github.com/web-auth-go/03_jwt/middleware"
)

func main() {

	initialisers.LoadEnv()
	initialisers.ConnectDB()
	initialisers.SyncDatabase()

	r := gin.Default()
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.GET("/validate", middleware.Auth, handlers.Validate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
