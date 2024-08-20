package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	lib "github.com/with-autro/autro-library"
	pb "github.com/with-autro/autro-price/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	host                string
	port                string
	serviceDiscoveryURL string
	priceService        pb.PriceServiceClient
)

func init() {
	host = os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port = os.Getenv("Port")
	if port == "" {
		port = "8080"
	}

	serviceDiscoveryURL = os.Getenv("SERVICE_DISCOVERY_URL")
	if serviceDiscoveryURL == "" {
		serviceDiscoveryURL = "http://autro-service-discovery:8500"
	}
}

// API Gateway 등록 함수
func registerService() {
	service := lib.Service{
		Name:    "autro-api-gateway",
		Address: fmt.Sprintf("%s:%s", host, port),
	}

	jsonData, err := json.Marshal(service)
	if err != nil {
		log.Fatalf("Failed to marshal service data: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/register", serviceDiscoveryURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to register service: %v. URL: %s", err, serviceDiscoveryURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to register service. Status code: %d", resp.StatusCode)
	}

	log.Println("Service registered successfully.")
}

// 서비스의 주소 가져오는 함수
func getServiceAddress(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/services/%s", serviceDiscoveryURL, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("service not found. Status Code: %d. Body: %s", resp.StatusCode, string(body))
	}

	var service lib.Service
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return "", err
	}
	return service.Address, nil
}

func initAutroPriceService() {
	for {
		addr, err := getServiceAddress("autro-price")
		if err != nil {
			log.Printf("Failed to get price service address: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to price service: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		priceService = pb.NewPriceServiceClient(conn)
		log.Println("Connected to price service")
		return
	}
}

// rest로 받으면 start gRPC를 쏘는 함수
func handleStart(c *gin.Context) {
	if priceService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Price service not available"})
		return
	}

	resp, err := priceService.Start(context.Background(), &pb.StartRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start price service: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func main() {
	// Service Discovery에 등록
	go registerService()

	// Init autro-price
	go initAutroPriceService()

	r := gin.Default()
	r.POST("/start", handleStart)

	log.Printf("Starting server on: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
