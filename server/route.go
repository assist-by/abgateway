package server

import (
	"github.com/gin-gonic/gin"
	"github.com/with-autro/autro-api-gateway/service"
)

func SetupRoutes(router *gin.Engine, serviceDiscoveryURL string) {
	router.POST("/start:autro-price", func(c *gin.Context) {
		service.StartPrice(c, serviceDiscoveryURL)
	})
}
