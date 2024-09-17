package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	lib "github.com/assist-by/autro-library"
	"github.com/segmentio/kafka-go"
)

// Service Discovery에 등록하는 함수
func RegisterService(writer *kafka.Writer, host, port string) error {
	service := lib.Service{
		Name:    "abgateway",
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
