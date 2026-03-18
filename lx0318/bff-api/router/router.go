package router

import (
	"github.com/gin-gonic/gin"
	"lx0318/bff-api/handler/service"
	"net/http"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	r.POST("/create", service.CreateOrder)
	return r
}
