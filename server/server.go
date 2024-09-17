package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/with-autro/autro-api-gateway/library"
)

type Server struct{
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *http.Server {
	router := gin.Default()
	SetupRoutes(router, cfg.ServiceDiscoveryURL)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}
}

func (s *)
