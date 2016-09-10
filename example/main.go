package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/ovh/go-ovh/ovh"
	"github.com/thbkrkr/go-utilz/http"
	"github.com/thbkrkr/ovh-auth"
)

var (
	buildDate = "dev"
	gitCommit = "dev"
)

func main() {
	flag.Parse()

	http.API("example", buildDate, gitCommit, router)
}

func router(r *gin.Engine) {
	authRouter := ovhauth.Secure(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ping"})
	})

	authRouter.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authRouter.GET("/api/me", func(c *gin.Context) {
		authUser := ovhauth.GetAuthUser(c)

		client, _ := ovh.NewDefaultClient()
		client.ConsumerKey = authUser.ConsumerKey

		var me map[string]string
		client.Get("/me", &me)

		c.JSON(200, me)
	})
}