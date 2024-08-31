package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	lib "github.com/with-autro/autro-library"
)

var (
	kafkaBroker         string
	host                string
	port                string
	registrationTopic   string
	serviceDiscoveryURL string
)

func init() {
	kafkaBroker = getEnv("KAFKA_BROKER", "kafka:9092")
	host = getEnv("HOST", "autro-api-gateway")
	port = getEnv("PORT", "50050")
	registrationTopic = getEnv("REGISTRATION_TOPIC", "service-registration")
	serviceDiscoveryURL = getEnv("SERVICE_DISCOVERY_URL", "http://autro-service-discovery:8500")

}

func getEnv(key, temp string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return temp
}

// Service Discovery에 등록하는 함수
func registerService(writer *kafka.Writer) error {
	service := lib.Service{
		Name:    "autro-api-gateway",
		Address: fmt.Sprintf("%s:%s", host, port),
	}

	jsonData, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("error marshaling service data: %v", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(service.Name),
		Value: jsonData,
	})

	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}

	log.Println("Service registration message sent successfully")
	return nil
}

// 서비스 등록 카프카 producer 생성
func createRegistrationWriter() *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:     []string{kafkaBroker},
			Topic:       registrationTopic,
			MaxAttempts: 5,
		})
}

// 서비스 주소 가져오는 API
func getServiceAddress(serviceName string) (string, error) {
	url :=
		fmt.Sprintf("%s/services/$s", serviceDiscoveryURL, serviceName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting service info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var service lib.Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling service data: %v", err)
	}

	return service.Address, nil
}

func startAutroPrice(c *gin.Context) {
	serviceAddress, err := getServiceAddress("autro-price")
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
	registrationWriter := createRegistrationWriter()
	defer registrationWriter.Close()

	if err := registerService(registrationWriter); err != nil {
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

}
