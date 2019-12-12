package main

import (
	"asterisk-dialer/api/originate"

	"github.com/gin-gonic/gin"
)

func initRoutes(r *gin.Engine) {
	// http://88.135.15.146:9002/api/originate/?number=521
	r.GET("/api/originate/", originate.Originate)
}
