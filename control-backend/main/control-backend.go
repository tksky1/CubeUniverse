package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello")
	})
	r.POST("/xxxpost", func())
	r.PUT("/xxxput", func())

	r.Run(":8080")
}
