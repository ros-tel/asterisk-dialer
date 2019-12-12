package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/heltonmarx/goami/ami"
)

func injectDB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
	}
}

func injectAmi(ami *ami.Socket) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ami", ami)
	}
}
