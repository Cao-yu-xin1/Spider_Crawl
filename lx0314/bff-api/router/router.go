package router

import (
	"github.com/gin-gonic/gin"
	"lx0314/bff-api/handler/service"
	"net/http"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	r.POST("/notify/pay", service.NotifyPay)
	return r
}
