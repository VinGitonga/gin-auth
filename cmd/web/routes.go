package main

import (
	"github.com/VinGitonga/gin-auth/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, g *handlers.GoApp) {
	router := r.Use(gin.Logger(), gin.Recovery())

	router.GET("/", g.Home())

	// set up  for storing details as cookies
	cookieData := cookie.NewStore([]byte("go-auth"))

	router.Use(sessions.Sessions("sessions", cookieData))

	router.POST("/sign-up", g.SignUp())
	router.POST("/sign-in", g.SignIn())

	authRouter := r.Group("/auth", Authorization())
	{
		authRouter.GET("/dashboard")
	}
}
