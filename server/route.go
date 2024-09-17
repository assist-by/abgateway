package server

import (
	"github.com/assist-by/abgateway/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, serviceDiscoveryURL string) {
	router.POST("/start:abprice", func(c *gin.Context) {
		service.StartPrice(c, serviceDiscoveryURL)
	})
}
