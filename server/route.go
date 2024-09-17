package server

import (
	"github.com/assist-by/autro-api-gateway/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, serviceDiscoveryURL string) {
	router.POST("/start:autro-price", func(c *gin.Context) {
		service.StartPrice(c, serviceDiscoveryURL)
	})
}
