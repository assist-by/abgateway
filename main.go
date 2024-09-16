package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/segmentio/kafka-go"
	config "github.com/with-autro/autro-api-gateway/library"
	kafka "github.com/with-autro/autro-api-gateway/pkg/kafka"
	service "github.com/with-autro/autro-api-gateway/service"
)

func startAutroPrice(c *gin.Context) {
	serviceAddress, err := service.GetServiceAddress("autro-price", cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get autro-price address: %v", err)})
		return
	}

	url := fmt.Sprintf("http://%s/start", serviceAddress)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start autro-price service : %v", err)})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response from autro-price service: %v", err)})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func main() {
	cfg := config.Load()

	registrationWriter := kafka.NewWriter(cfg.KafkaBroker, cfg.RegistrationTopic)
	defer registrationWriter.Close()

	if err := service.RegisterService(registrationWriter, cfg.Host, cfg.Port); err != nil {
		log.Printf("Failed to register service: %v\n", err)
	}

	router := gin.Default()
	router.POST("/start:autro-price", startAutroPrice)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
