package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartPrice(c *gin.Context, serviceDiscoveryURL string) {
	serviceAddress, err := GetServiceAddress("autro-price", serviceDiscoveryURL)
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
